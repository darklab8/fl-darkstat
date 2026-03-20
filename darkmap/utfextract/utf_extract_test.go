package utfextract

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

// buildUTF constructs a minimal valid UTF binary containing one texture
// library entry with a single image leaf (MIP0).
//
// Tree layout:
//
//	node0 @ tree[0]:  '\'           DIR  sib=0    child=node1
//	node1 @ tree[44]: 'Texture library' DIR  sib=0    child=node2
//	node2 @ tree[88]: <imageName>   DIR  sib=0    child=node3
//	node3 @ tree[132]:'MIP0'        LEAF sib=0    child=0 (data offset 0)
func buildUTF(t *testing.T, imageName string, imageData []byte) []byte {
	t.Helper()

	// ── String segment ────────────────────────────────────────────────────
	var strBuf bytes.Buffer
	offsets := map[string]uint32{}
	addStr := func(s string) uint32 {
		if o, ok := offsets[s]; ok {
			return o
		}
		o := uint32(strBuf.Len())
		strBuf.WriteString(s)
		strBuf.WriteByte(0)
		offsets[s] = o
		return o
	}
	oBackslash := addStr("\\")
	oTexLib := addStr("Texture library")
	oImgEntry := addStr(imageName)
	oMip0 := addStr("MIP0")
	strBytes := strBuf.Bytes()

	// ── Tree nodes ────────────────────────────────────────────────────────
	packNode := func(n utfNode) []byte {
		buf := make([]byte, 44)
		binary.LittleEndian.PutUint32(buf[0:], n.SiblingOffset)
		binary.LittleEndian.PutUint32(buf[4:], n.NameOffset)
		binary.LittleEndian.PutUint32(buf[8:], n.Flags)
		binary.LittleEndian.PutUint32(buf[12:], n.Zero)
		binary.LittleEndian.PutUint32(buf[16:], n.ChildOffset)
		binary.LittleEndian.PutUint32(buf[20:], n.AllocSize)
		binary.LittleEndian.PutUint32(buf[24:], n.Size1)
		binary.LittleEndian.PutUint32(buf[28:], n.Size2)
		return buf
	}

	sz := uint32(len(imageData))
	var treeBuf bytes.Buffer
	// node0 @ tree[0]:   '\' directory, child @ tree[44]
	treeBuf.Write(packNode(utfNode{NameOffset: oBackslash, Flags: flagHasChildren, ChildOffset: 44}))
	// node1 @ tree[44]:  'Texture library' directory, child @ tree[88]
	treeBuf.Write(packNode(utfNode{NameOffset: oTexLib, Flags: flagHasChildren, ChildOffset: 88}))
	// node2 @ tree[88]:  imageName directory, child @ tree[132]
	treeBuf.Write(packNode(utfNode{NameOffset: oImgEntry, Flags: flagHasChildren, ChildOffset: 132}))
	// node3 @ tree[132]: 'MIP0' leaf, data at data[0]
	treeBuf.Write(packNode(utfNode{NameOffset: oMip0, Flags: flagLeaf, ChildOffset: 0, Size1: sz, Size2: sz, AllocSize: sz}))
	treeBytes := treeBuf.Bytes()

	// ── Assemble file ─────────────────────────────────────────────────────
	// "UTF " (4) + header (40) + 12-byte pad + tree + strings + data
	pad := []byte("000000000000")
	treeOffset := uint32(4 + 40 + len(pad))
	strOffset := treeOffset + uint32(len(treeBytes))
	dataOffset := strOffset + uint32(len(strBytes))

	var hdr bytes.Buffer
	for _, v := range []uint32{
		0x101,
		treeOffset, uint32(len(treeBytes)),
		0, 44, // treefirst, treeelemsize
		strOffset, uint32(len(strBytes)), uint32(len(strBytes)),
		dataOffset, 0, // datafirst
	} {
		binary.Write(&hdr, binary.LittleEndian, v)
	}

	var out bytes.Buffer
	out.WriteString("UTF ")
	out.Write(hdr.Bytes())
	out.Write(pad)
	out.Write(treeBytes)
	out.Write(strBytes)
	out.Write(imageData)
	return out.Bytes()
}

// ─────────────────────────────────────────────────────────────────────────────

func TestParseBadMagic(t *testing.T) {
	_, err := ParseUTF([]byte("NOPE"))
	if err == nil {
		t.Fatal("expected error for bad magic, got nil")
	}
}

func TestRootNodeAtOffsetZeroIsEntered(t *testing.T) {
	// The root directory node lives at tree offset 0.
	// A previous bug used `for offset != 0` which skipped it entirely.
	fakeImg := make([]byte, 32) // TGA (non-DDS)
	raw := buildUTF(t, "test.tga", fakeImg)

	root, err := ParseUTF(raw)
	if err != nil {
		t.Fatalf("ParseUTF failed: %v", err)
	}
	if len(root.Children) == 0 {
		t.Fatal("root has no children: offset-0 node was not entered")
	}
	backslash := root.Child("\\")
	if backslash == nil {
		t.Fatalf("'\\' node not found; got children: %v", childNames(root))
	}
	texLib := backslash.Child("Texture library")
	if texLib == nil {
		t.Fatalf("'Texture library' not found under '\\'; got: %v", childNames(backslash))
	}
}

func TestExtractTGA(t *testing.T) {
	fakeImg := make([]byte, 32) // starts with 0x00 → not DDS
	raw := buildUTF(t, "mytex.tga", fakeImg)

	images, err := ExtractImages(raw)
	if err != nil {
		t.Fatalf("ExtractImages failed: %v", err)
	}
	if len(images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(images))
	}
	if images[0].Filename != "mytex.tga" {
		t.Errorf("expected mytex.tga, got %s", images[0].Filename)
	}
	if !bytes.Equal(images[0].Data, fakeImg) {
		t.Error("image data mismatch")
	}
}

func TestExtractDDS(t *testing.T) {
	// File named .tga but contains DDS data → extension must be rewritten
	fakeImg := []byte{'D', 'D', 'S', ' ', 0x7C, 0x00, 0x00, 0x00}
	raw := buildUTF(t, "mytex.tga", fakeImg)

	images, err := ExtractImages(raw)
	if err != nil {
		t.Fatalf("ExtractImages failed: %v", err)
	}
	if len(images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(images))
	}
	if images[0].Filename != "mytex.dds" {
		t.Errorf("expected mytex.dds (DDS detected), got %s", images[0].Filename)
	}
}

func TestFindTextureLibraryCaseInsensitive(t *testing.T) {
	root := &Node{Name: "root", Children: []*Node{
		{Name: "TEXTURE LIBRARY", Children: []*Node{}},
	}}
	if findTextureLibrary(root) == nil {
		t.Fatal("findTextureLibrary should be case-insensitive")
	}
}

func TestIsDDS(t *testing.T) {
	cases := []struct {
		data []byte
		want bool
	}{
		{[]byte{'D', 'D', 'S'}, true},
		{[]byte{'D', 'D', 'S', 0x20}, true},
		{[]byte{0x00, 0x00, 0x00}, false},
		{[]byte{}, false},
	}
	for _, c := range cases {
		if got := isDDS(c.data); got != c.want {
			t.Errorf("isDDS(%v) = %v, want %v", c.data, got, c.want)
		}
	}
}

func TestReplaceTGAExt(t *testing.T) {
	cases := []struct{ in, want string }{
		{"foo.tga", "foo.dds"},
		{"foo.TGA", "foo.dds"},
		{"foo.png", "foo.png.dds"},
		{"noext", "noext.dds"},
	}
	for _, c := range cases {
		if got := replaceTGAExt(c.in, ".dds"); got != c.want {
			t.Errorf("replaceTGAExt(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestIsUTFExtension(t *testing.T) {
	for _, f := range []string{"model.3db", "ship.CMP", "tex.txm", "mat.MAT", "dfm.dfm"} {
		if !isUTFExtension(f) {
			t.Errorf("isUTFExtension(%q) should be true", f)
		}
	}
	for _, f := range []string{"readme.txt", "image.png", "file.exe"} {
		if isUTFExtension(f) {
			t.Errorf("isUTFExtension(%q) should be false", f)
		}
	}
}

// childNames returns the names of n's direct children (for test diagnostics).
func childNames(n *Node) []string {
	names := make([]string, len(n.Children))
	for i, c := range n.Children {
		names[i] = c.Name
	}
	return names
}

// ─────────────────────────────────────────────────────────────────────────────
// Path-preservation tests
// ─────────────────────────────────────────────────────────────────────────────

// TestExtractFromDirPreservePaths verifies that with preservePaths=true the
// sub-directory structure of the input root is mirrored under outputDir.
//
// Input layout:
//
//	<tmpIn>/
//	  a/
//	    b/
//	      tex.txm   (one TGA image inside)
//
// Expected output (preservePaths=true, recursive=true):
//
//	<tmpOut>/a/b/tex.txm/<image>.tga
//
// Expected output (preservePaths=false, recursive=true):
//
//	<tmpOut>/tex.txm/<image>.tga   ← sub-dirs stripped
func TestExtractFromDirPreservePaths(t *testing.T) {
	fakeImg := make([]byte, 32) // TGA
	utfData := buildUTF(t, "icon.tga", fakeImg)

	// Build  <tmpIn>/a/b/tex.txm
	tmpIn := t.TempDir()
	subDir := filepath.Join(tmpIn, "a", "b")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}
	utfPath := filepath.Join(subDir, "tex.txm")
	if err := os.WriteFile(utfPath, utfData, 0o644); err != nil {
		t.Fatal(err)
	}

	t.Run("preserve=true", func(t *testing.T) {
		tmpOut := t.TempDir()
		fr, iw, err := ExtractFromDir(tmpIn, tmpOut, true, true)
		if err != nil {
			t.Fatalf("ExtractFromDir: %v", err)
		}
		if fr != 1 || iw != 1 {
			t.Fatalf("expected 1 file / 1 image, got files=%d images=%d", fr, iw)
		}
		want := filepath.Join(tmpOut, "a", "b", "tex.txm", "icon.tga")
		if _, err := os.Stat(want); os.IsNotExist(err) {
			t.Errorf("expected output at %s (not found)", want)
		}
	})

	t.Run("preserve=false", func(t *testing.T) {
		tmpOut := t.TempDir()
		fr, iw, err := ExtractFromDir(tmpIn, tmpOut, true, false)
		if err != nil {
			t.Fatalf("ExtractFromDir: %v", err)
		}
		if fr != 1 || iw != 1 {
			t.Fatalf("expected 1 file / 1 image, got files=%d images=%d", fr, iw)
		}
		// Sub-dirs NOT preserved — file sits directly under tex.txm/
		want := filepath.Join(tmpOut, "tex.txm", "icon.tga")
		if _, err := os.Stat(want); os.IsNotExist(err) {
			t.Errorf("expected flat output at %s (not found)", want)
		}
		// And the deep path must NOT exist
		deep := filepath.Join(tmpOut, "a", "b", "tex.txm", "icon.tga")
		if _, err := os.Stat(deep); err == nil {
			t.Errorf("deep path %s should not exist when preservePaths=false", deep)
		}
	})
}
