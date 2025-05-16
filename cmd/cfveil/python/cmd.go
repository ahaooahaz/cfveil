package python

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/ahaooahaz/cfveil/internal/python"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "python",
	Short:   "python project.",
	Long:    `python project.`,
	Aliases: []string{"py"},
	Run: func(cmd *cobra.Command, args []string) {
		if *arg_INPUT == "" || *arg_OUTPUT == "" {
			fmt.Println("invalid")
			return
		}

		err := python.Process(*arg_INPUT, *arg_OUTPUT, *arg_EXCLUDE)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		return

		// xs := []*X{}
		// for k, v := range pyFiles {
		// 	x := &X{}
		// 	x.Src = v

		// 	{
		// 		info, _ := os.Stat(v)
		// 		suffix := ""
		// 		if !info.IsDir() {
		// 			suffix = ".py"
		// 		}
		// 		s := k
		// 		idx := strings.LastIndex(s, ".")
		// 		if idx != -1 {
		// 			s = s[:idx]
		// 		}

		// 		x.ImportSrc = strings.ReplaceAll(s, "/", ".")
		// 		parts := strings.Split(x.ImportSrc, ".")
		// 		for i := 0; i < len(parts); i++ {
		// 			hash := md5.Sum([]byte(parts[i]))
		// 			parts[i] = "x" + hex.EncodeToString(hash[:])
		// 		}
		// 		x.ImportDst = strings.Join(parts, ".")

		// 		parts = strings.Split(s, "/")
		// 		for i := 0; i < len(parts); i++ {
		// 			hash := md5.Sum([]byte(parts[i]))
		// 			parts[i] = "x" + hex.EncodeToString(hash[:])
		// 		}
		// 		x.Dst = strings.Join(parts, "/")
		// 		x.Dst += suffix
		// 	}
		// 	xs = append(xs, x)
		// }

		// for _, x := range xs {
		// 	dstfile := "dist/" + x.Dst
		// 	if err := os.MkdirAll(filepath.Dir(dstfile), os.ModePerm); err != nil {
		// 		panic(err)
		// 	}
		// 	CopyFile(x.Src, dstfile)
		// 	processFile(dstfile, xs)
		// }
	},
}

func processFile(path string, xs []*X) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// 编译正则：匹配 from xxx 或 import xxx 开头
		re := regexp.MustCompile(`^\s*(from|import)\b.*`)

		if re.MatchString(strings.TrimSpace(line)) {
			for _, xx := range xs {
				line = strings.ReplaceAll(line, xx.ImportSrc, xx.ImportDst)
			}
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// 将结果写回原文件（你也可以写到另一个文件）
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}

func CopyFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 获取原文件权限
	fileInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fileInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

type X struct {
	Src       string
	Dst       string
	ImportSrc string
	ImportDst string
}

var (
	arg_INPUT, arg_OUTPUT *string
	arg_EXCLUDE           *[]string
)

func init() {
	arg_INPUT = Cmd.Flags().StringP("input", "i", "", "project dir or file path")
	arg_OUTPUT = Cmd.Flags().StringP("output", "o", "", "output dir")
	arg_EXCLUDE = Cmd.Flags().StringSliceP("exclude", "e", []string{}, "exclude files")
}
