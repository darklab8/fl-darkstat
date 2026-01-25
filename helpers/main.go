package helpers

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/exe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/infocard_mapped/infocard"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/cantil"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func HelpersCliGroup(Args []string) {

	fmt.Println("freelancer folder=", settings.Env.FreelancerFolder, settings.Env)
	parser := cantil.NewConsoleParser(
		[]cantil.Action{
			{
				Nickname:    "infocard_read",
				Description: "Read some infocard and output to file",
				Func: func(info cantil.ActionInfo) error {
					fmt.Println("inputed args=", info.CmdArgs[1:])
					workdir, _ := os.Getwd()

					os.Args = info.CmdArgs

					var input_file string
					flag.StringVar(&input_file, "input_file", "", "Name of file input")

					var output_file string
					flag.StringVar(&output_file, "out_file", "", "Name of file output")
					var output_stdout bool
					flag.BoolVar(&output_stdout, "out_stdout", false, "Specify to output to stdout")
					var dll_index_num int
					index_error := `Missing argument "--int_index". Index files begin with specific global_offset global_offset := int(math.Pow(2, 16)) * (idx)
You need to spcify index number of a file you are reading. it is expected integer from 0 to 10 for example, it matches row number of DLL record in [Resources] in freelancer.ini file
[Resources]
DLL = InfoCards.dll ; matches number 0
DLL = MiscText.dll ; matches number 1
DLL = NameResources.dll	; matches number 2
DLL = EquipResources.dll 	; matches number 3
DLL = OfferBribeResources.dll	; matches number 4
DLL = MiscTextInfo2.dll	; matches number 5
DLL = Discovery.dll	; matches number 6
DLL = DsyAddition.dll	; matches number 7
`
					flag.IntVar(&dll_index_num, "int_index", -1, index_error)

					flag.Parse()
					args := flag.Args()

					if len(args) > 0 {
						logus.Log.Error("found not supported not parsed args", typelog.Any("args", args))
						os.Exit(1)
					}

					if input_file == "" {
						logus.Log.Error("Missing argument argument of a file input . Expected to get as input filename of infocard. Input is case sensetive.", typelog.Any("args", flag.Args()))
						os.Exit(1)
					}
					filename := input_file

					if dll_index_num == -1 {
						logus.Log.Error(index_error)
						os.Exit(1)
					}

					if !output_stdout && output_file == "" {
						logus.Log.Error("specify out_stdout=true or out_file=path_to_file for output")
						os.Exit(1)
					}

					file_path := utils_types.FilePath(workdir).Join(filename)
					data, err := os.ReadFile(file_path.ToString())

					if logus.Log.CheckError(err, "unable to read dll") {
						os.Exit(1)
					}

					// if you inject "resources.dll" as 0 element of the list to process
					// despite it being not present in freelancer.ini and original Alex parsing script
					// then we go with global_offset from (idx) instead of (idx+1) as Alex had.
					global_offset := int(math.Pow(2, 16)) * (dll_index_num)

					out := infocard.NewConfig()
					exe_mapped.ParseDLL(data, out, global_offset)
					config := out.GetUnsafe()

					if len(config.Infocards) == 0 && len(config.Infonames) == 0 {
						logus.Log.Error("not found any infocards inside")
						os.Exit(1)
					}

					type Infoname struct {
						num  int
						card infocard.Infoname
					}
					type Infocard struct {
						num  int
						card *infocard.Infocard
					}
					var infonames []Infoname
					var infocards []Infocard
					for name_id, content := range config.Infonames {
						infonames = append(infonames, Infoname{
							num:  name_id,
							card: content,
						})
					}
					for name_id, content := range config.Infocards {
						infocards = append(infocards, Infocard{
							num:  name_id,
							card: content,
						})
					}
					sort.Slice(infonames[:], func(i, j int) bool {
						return infonames[i].num < infonames[j].num
					})
					sort.Slice(infocards[:], func(i, j int) bool {
						return infocards[i].num < infocards[j].num
					})

					if output_stdout {
						for _, content := range infonames {
							fmt.Printf("[%d]\n", content.num)
							fmt.Println(content.card)
						}

						for _, content := range infocards {
							fmt.Printf("[%d]\n", content.num)
							fmt.Println(content.card.GetContent())
						}
					}

					if output_file != "" {
						filepath := filepath.Join(workdir, "output_file")
						f, err := os.Create(filepath)
						logus.Log.CheckPanic(err, "failed to open file for reading")

						for name_id, content := range config.Infonames {
							_, err = fmt.Fprintf(f, "[%d]\n", name_id)
							logus.Log.CheckPanic(err, "failed to write to file")
							_, err = fmt.Fprintln(f, string(content))
							logus.Log.CheckPanic(err, "failed to write to file")
						}

						for name_id, content := range config.Infocards {
							_, err = fmt.Fprintf(f, "[%d]\n", name_id)
							logus.Log.CheckPanic(err, "failed to write to file")
							_, err = fmt.Fprintln(f, string(content.GetContent()))
							logus.Log.CheckPanic(err, "failed to write to file")
						}
					}

					return nil
				},
			},
		},
		cantil.ParserOpts{
			ParentArgs: []string{"helpers"},
			Enverants:  settings.Enverants,
		},
	)
	err := parser.Run(Args)
	logus.Log.CheckError(err, "failed to execute helpers cli group command")
}
