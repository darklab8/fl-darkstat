// Package utfextract provides parsing and image extraction for
// Microsoft Freelancer (2003) UTF binary files (.txm, .cmp, .3db, .mat, .dfm).
//
// The UTF format is a tree of named nodes, each node holding either
// child nodes or raw binary data.  Texture images (TGA or DDS) are stored
// inside a sub-tree called "Texture library".
package utfextract

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

// ────────────────────────────────────────────────────────────────────────────
// Low-level UTF on-disk structures
// ────────────────────────────────────────────────────────────────────────────

// utfHeader is the 40-byte block that follows the 4-byte "UTF " magic.
type utfHeader struct {
	Version      uint32
	TreeOffset   uint32
	TreeSize     uint32
	TreeFirst    uint32 // reserved / unused
	TreeElemSize uint32 // always 44 = sizeof(utfNode)
	StringOffset uint32
	StringSpace  uint32
	StringSize   uint32
	DataOffset   uint32
	DataFirst    uint32 // reserved / unused
}

// utfNode is the 44-byte record for one entry in the tree segment.
type utfNode struct {
	SiblingOffset uint32
	NameOffset    uint32
	Flags         uint32
	Zero          uint32
	ChildOffset   uint32
	AllocSize     uint32
	Size1         uint32
	Size2         uint32
	Time1         uint32
	Time2         uint32
	Time3         uint32
}

const (
	flagHasChildren = 0x10 // node is a directory; ChildOffset points into tree
	flagLeaf        = 0x80 // node is a leaf;      ChildOffset points into data

	nodeSizeBytes = 44
)

// ────────────────────────────────────────────────────────────────────────────
// Parsed tree node
// ────────────────────────────────────────────────────────────────────────────

// Node represents one entry in the UTF tree.
// Exactly one of Children or Data is non-nil.
type Node struct {
	Name     string
	Children []*Node // directory node
	Data     []byte  // leaf node
}

// Child returns the first direct child whose name matches (case-insensitive).
func (n *Node) Child(name string) *Node {
	for _, c := range n.Children {
		if strings.EqualFold(c.Name, name) {
			return c
		}
	}
	return nil
}

// ────────────────────────────────────────────────────────────────────────────
// Parser
// ────────────────────────────────────────────────────────────────────────────

type parser struct {
	tree    []byte
	strings []byte
	data    []byte
	visited map[uint32]bool
}

// ParseUTF reads raw UTF file bytes and returns the root node of the tree.
// The root node is a synthetic container whose Children are the top-level
// entries of the file (typically a single "\" directory node).
func ParseUTF(raw []byte) (*Node, error) {
	if len(raw) < 4 || string(raw[:4]) != "UTF " {
		return nil, errors.New("not a UTF file (bad magic)")
	}
	if len(raw) < 44 { // 4 magic + 40 header bytes
		return nil, errors.New("UTF file too short to contain header")
	}

	var hdr utfHeader
	if err := binary.Read(bytes.NewReader(raw[4:44]), binary.LittleEndian, &hdr); err != nil {
		return nil, fmt.Errorf("reading UTF header: %w", err)
	}

	treeEnd := int(hdr.TreeOffset) + int(hdr.TreeSize)
	strEnd := int(hdr.StringOffset) + int(hdr.StringSpace)
	if treeEnd > len(raw) || strEnd > len(raw) || int(hdr.DataOffset) > len(raw) {
		return nil, errors.New("UTF segment extends past end of file")
	}

	p := &parser{
		tree:    raw[hdr.TreeOffset:treeEnd],
		strings: raw[hdr.StringOffset:strEnd],
		data:    raw[hdr.DataOffset:],
		visited: make(map[uint32]bool),
	}

	root := &Node{Name: "<root>"}
	// The tree segment begins with the root entry at offset 0.
	// parseNode unconditionally enters offset 0, then follows the sibling
	// chain -- matching Perl's UTFreadUTFrek($tree, 0) behaviour.
	p.parseNode(0, root)
	return root, nil
}

// parseNode parses the entry at treeOffset, appends it to parent.Children,
// then recursively follows the sibling chain (also appended to parent.Children).
//
// KEY INVARIANT: treeOffset == 0 is the *root node*, not a null pointer.
// A null/end-of-chain is signalled only by SiblingOffset == 0 coming out of
// an already-parsed node's sibling field -- never by the initial call argument.
// The old code used `for offset != 0 { ... }` which skipped offset 0 entirely.
func (p *parser) parseNode(treeOffset uint32, parent *Node) {
	if p.visited[treeOffset] {
		return
	}
	if int(treeOffset)+nodeSizeBytes > len(p.tree) {
		return
	}
	p.visited[treeOffset] = true

	var rec utfNode
	if err := binary.Read(
		bytes.NewReader(p.tree[treeOffset:treeOffset+nodeSizeBytes]),
		binary.LittleEndian, &rec,
	); err != nil {
		return
	}

	name := p.readString(rec.NameOffset)

	// Effective size: smaller of Size1/Size2, but ignore Size2 when zero
	// (zero means "not set", not "zero bytes").
	size := rec.Size1
	if rec.Size2 != 0 && rec.Size2 < size {
		size = rec.Size2
	}

	node := &Node{Name: name}

	isDir := (rec.Flags&flagHasChildren) != 0 && (rec.Flags&flagLeaf) == 0
	if isDir {
		node.Children = []*Node{}
		if rec.ChildOffset != 0 {
			p.parseNode(rec.ChildOffset, node)
		}
	} else {
		node.Data = p.readData(rec.ChildOffset, size)
	}

	parent.Children = append(parent.Children, node)

	// SiblingOffset == 0 means end of chain.
	if rec.SiblingOffset != 0 {
		p.parseNode(rec.SiblingOffset, parent)
	}
}

func (p *parser) readString(offset uint32) string {
	if int(offset) >= len(p.strings) {
		return ""
	}
	rest := p.strings[offset:]
	end := bytes.IndexByte(rest, 0)
	if end < 0 {
		return string(rest)
	}
	return string(rest[:end])
}

func (p *parser) readData(offset, size uint32) []byte {
	if size == 0 || int(offset) >= len(p.data) {
		return nil
	}
	end := int(offset) + int(size)
	if end > len(p.data) {
		end = len(p.data)
	}
	out := make([]byte, end-int(offset))
	copy(out, p.data[int(offset):end])
	return out
}

// ────────────────────────────────────────────────────────────────────────────
// Image extraction
// ────────────────────────────────────────────────────────────────────────────

// ExtractedImage holds one image found inside a UTF file.
type ExtractedImage struct {
	// Filename is the output path relative to the UTF file's output directory.
	// May contain one sub-directory for multi-mip entries, e.g. "asterdust/MIP0.tga".
	Filename string
	Data     []byte
}

// ExtractImages parses a raw UTF file and returns all texture images found
// in its "Texture library" sub-tree.
func ExtractImages(raw []byte) ([]ExtractedImage, error) {
	root, err := ParseUTF(raw)
	if err != nil {
		return nil, err
	}

	// Normal tree layout:  root -> '\' -> 'Texture library' -> <entries>
	// findTextureLibrary searches the immediate children of its argument,
	// so we try root directly first, then the '\' child.
	texLib := findTextureLibrary(root)
	if texLib == nil {
		if bs := root.Child("\\"); bs != nil {
			texLib = findTextureLibrary(bs)
		}
	}
	if texLib == nil {
		return nil, errors.New("no texture library found in file")
	}

	var images []ExtractedImage

	for _, texEntry := range texLib.Children {
		if len(texEntry.Children) == 0 {
			continue
		}

		if len(texEntry.Children) == 1 {
			// Single child (one MIP0 or MIPS): save as <TextureName>.tga/.dds
			child := texEntry.Children[0]
			if len(child.Data) == 0 {
				continue
			}
			fname := texEntry.Name
			if isDDS(child.Data) {
				fname = replaceTGAExt(fname, ".dds")
			} else {
				fname = ensureExt(fname, ".tga")
			}
			images = append(images, ExtractedImage{Filename: fname, Data: child.Data})
			continue
		}

		// Multiple children (MIP0..MIPn, or .dfm sub-entries):
		// save as <TextureName>/<MIPKey>.tga (or .dds)
		for _, child := range texEntry.Children {
			if len(child.Data) == 0 {
				continue
			}
			ext := ".tga"
			if isDDS(child.Data) {
				ext = ".dds"
			}
			fname := texEntry.Name + "/" + child.Name + ext
			images = append(images, ExtractedImage{Filename: fname, Data: child.Data})
		}
	}

	return images, nil
}

func findTextureLibrary(parent *Node) *Node {
	for _, c := range parent.Children {
		if strings.EqualFold(strings.TrimSpace(c.Name), "texture library") {
			return c
		}
	}
	return nil
}

func isDDS(data []byte) bool {
	return len(data) >= 3 && data[0] == 'D' && data[1] == 'D' && data[2] == 'S'
}

func replaceTGAExt(name, newExt string) string {
	if strings.HasSuffix(strings.ToLower(name), ".tga") {
		return name[:len(name)-4] + newExt
	}
	return name + newExt
}

func ensureExt(name, ext string) string {
	if filepath.Ext(name) == "" {
		return name + ext
	}
	return name
}

// ────────────────────────────────────────────────────────────────────────────
// High-level file helpers
// ────────────────────────────────────────────────────────────────────────────

// ExtractFromFile reads one UTF file and writes all found images into
//
//	outputDir/<basename>/<imagefile>
//
// Returns the number of images written.
func ExtractFromFile(inputPath, outputDir string, shapes *Shapes) error {
	return extractFromFileWithSubdir(inputPath, outputDir, "", shapes)
}

type ShapeOption func(s *Shapes)

func NewShapes(opts ...ShapeOption) *Shapes {
	s := &Shapes{
		ShapesByNick: make(map[string]*Shape),
		ShapesConfig: ShapesConfig{
			PermittedShapes: make(map[string]bool),
		},
	}

	for _, opt := range opts {
		opt(s)
	}
	return s
}

// extractFromFileWithSubdir is the internal implementation that also accepts
// an optional subdir (relative path of the UTF file's parent directory
// measured from the original input root).  When non-empty it is inserted
// between outputDir and the UTF file's basename so that the full output path
// becomes:
//
//	outputDir/<subdir>/<basename>/<imagefile>
//
// This preserves the original directory tree when extracting recursively,
// making it easy to diff output against a reference tree.
func extractFromFileWithSubdir(inputPath, outputDir, subdir string, shapes *Shapes) error {
	raw, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", inputPath, err)
	}

	images, err := ExtractImages(raw)
	if err != nil {
		return fmt.Errorf("extracting from %s: %w", inputPath, err)
	}

	basename := filepath.Base(inputPath)
	for _, img := range images {
		dest := filepath.Join(outputDir, subdir, basename, filepath.FromSlash(img.Filename))
		if strings.Contains(img.Filename, semantic.PATH_SEPARATOR) {
			basename = filepath.Dir(img.Filename)
		}
		img.Filename = filepath.Base(img.Filename)
		name := strings.Split(strings.ToLower(img.Filename), ".")

		// if strings.Contains(strings.ToLower(dest), "navmaptextures.txm") {
		// 	fmt.Print()
		// }

		image := &Image{
			Nickname:  name[0],
			Extension: name[1],
			Data:      img.Data,
			Dest:      dest,
		}

		var shape_nickname, shape_extension string
		if strings.Contains(basename, ".") {
			shape_naming := strings.Split(strings.ToLower(basename), ".")
			shape_nickname = shape_naming[0]
			shape_extension = shape_naming[1]
		} else {
			shape_nickname = basename
		}

		shape, ok := shapes.ShapesByNick[shape_nickname]
		if !ok {
			shape = &Shape{
				Nickname:  shape_nickname,
				Extension: shape_extension,
				Path:      dest,
			}
			shapes.ShapesByNick[shape_nickname] = shape
		}
		shape.Images = append(shape.Images, image)

		if shapes.WriteToFile {
			if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
				return fmt.Errorf("creating directory for %s: %w", dest, err)
			}
			if err := os.WriteFile(dest, img.Data, 0o644); err != nil {
				return fmt.Errorf("writing %s: %w", dest, err)
			}
		}

		shapes.ImageWritten++
	}

	return nil
}

type ShapesConfig struct {
	WriteToFile bool

	//optimization to use less for purpose of in memory faster dev run
	UsePermittedShapesWhitelist bool
	PermittedShapes             map[string]bool
}
type Shapes struct {
	ShapesConfig
	ImageWritten int
	FilesRead    int
	ShapesByNick map[string]*Shape
}

type Shape struct {
	Nickname  string
	Extension string
	Images    []*Image
	Path      string
}
type Image struct {
	Nickname  string
	Extension string
	Dest      string
	Data      []byte
}

// ExtractFromDir walks inputDir for UTF files and extracts images into outputDir.
// When recursive is true it descends into sub-directories.
//
// When preservePaths is true the sub-directory structure inside inputDir is
// mirrored under outputDir, so a file found at
//
//	inputDir/foo/bar/model.3db
//
// produces images under
//
//	outputDir/foo/bar/model.3db/<imagefile>
//
// rather than the flat
//
//	outputDir/model.3db/<imagefile>
//
// Use preservePaths=true when validating output against a known-good reference
// tree extracted with the Perl tool.
func ExtractFromDir(inputDir, outputDir string, recursive, preservePaths bool, shapes *Shapes) (err error) {
	return extractDir(inputDir, outputDir, "", recursive, preservePaths, shapes)
}

// extractDir is the recursive implementation; subdir accumulates the relative
// path from the original inputDir root as we descend.
func extractDir(inputDir, outputDir, subdir string, recursive, preservePaths bool, shapes *Shapes) (err error) {
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("reading directory %s: %w", inputDir, err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(inputDir, entry.Name())

		if entry.IsDir() {
			if recursive {
				childSubdir := entry.Name()
				if subdir != "" {
					childSubdir = filepath.Join(subdir, entry.Name())
				}
				e := extractDir(fullPath, outputDir, childSubdir, recursive, preservePaths, shapes)

				if e != nil {
					err = e
				}
			}
			continue
		}

		if !isUTFExtension(entry.Name()) {
			continue
		}

		shapes.FilesRead++
		outSubdir := ""
		if preservePaths {
			outSubdir = subdir
		}
		e := extractFromFileWithSubdir(fullPath, outputDir, outSubdir, shapes)
		if e != nil {
			fmt.Fprintf(os.Stderr, "warning: %v\n", e)
		}
	}

	return err
}

func isUTFExtension(fname string) bool {
	switch strings.ToLower(filepath.Ext(fname)) {
	case ".3db", ".cmp", ".mat", ".txm", ".dfm":
		return true
	}
	return false
}
