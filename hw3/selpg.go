package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/pflag"
)

//a sruct for selpg_args
type selpgArgs struct {
	startPage  int
	endPage    int
	inFilename string
	pageLen    int
	pageType   bool

	printDest string
}

var progname string

//INTMAX is the maxinum value 'signed int' can hold
const INTMAX = int(^uint(0) >> 1)

//Usage is used for telling the use of selpg
func Usage() {
	fmt.Fprintf(os.Stderr, "\nUSAGE: %s -sstartPage -eendPage [ -f | -llinesPerPage ] [ -ddest ] [ inFilename ]\n", progname)
}

//ProcessArgs process the args
func ProcessArgs(sa *selpgArgs) {
	//pflag 绑定参数
	pflag.IntVarP(&sa.startPage, "startPage", "s", 1, "Start page number")
	pflag.IntVarP(&sa.endPage, "endPage", "e", 1, "End page number")
	pflag.IntVarP(&sa.pageLen, "pageLen", "l", 7, "Lines per page")
	pflag.BoolVarP(&sa.pageType, "pageType", "f", false, "Page type")
	pflag.StringVarP(&sa.printDest, "dest", "d", "", "Destination")
	pflag.Usage = func() {
		Usage()
		pflag.PrintDefaults()
	}
	pflag.Parse()

	sa.inFilename = ""
	if remain := pflag.Args(); len(remain) > 0 {
		sa.inFilename = remain[0]
	}

	//check the args are valid
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s: not enough arguments\n", progname)
		pflag.Usage()
		os.Exit(1)
	}
	if sa.startPage < 1 || sa.startPage > (INTMAX-1) {
		fmt.Fprintf(os.Stderr, "%s: invalid start page %d\n", progname, sa.startPage)
		pflag.Usage()
		os.Exit(2)
	}
	if sa.endPage < 1 || sa.endPage > (INTMAX-1) || sa.endPage < sa.startPage {
		fmt.Fprintf(os.Stderr, "%s: invalid end page %d\n", progname, sa.endPage)
		pflag.Usage()
		os.Exit(3)
	}
	if sa.pageLen < 1 || sa.pageLen > (INTMAX-1) {
		fmt.Fprintf(os.Stderr, "%s: invalid page length %d\n", progname, sa.pageLen)
		pflag.Usage()
		os.Exit(4)
	}
	if sa.inFilename != "" {
		if _, err := os.Stat(sa.inFilename); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s: input file \"%s\" does not exist\n", progname, sa.inFilename)
			pflag.Usage()
			os.Exit(5)
		}
	}
}

//ProcessInput complete the program running method
func ProcessInput(sa *selpgArgs) {
	filein := os.Stdin
	fileout := os.Stdout
	lineCount := 0
	pageCount := 1

	//读取文件并判断读取是否出错
	if sa.inFilename != "" {
		err := errors.New("")
		filein, err = os.Open(sa.inFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: could not open input file \"%s\"\n", progname, sa.inFilename)
			os.Exit(6)
		}
		defer filein.Close()
	}

	//开始根据是否以换页符分页进行分页
	readLine := bufio.NewReader(filein)
	if sa.pageType == false {
		for {
			line, err := readLine.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: Read file error!\n", progname)
				os.Exit(7)
			}
			lineCount++
			if lineCount > sa.pageLen {
				pageCount++
				lineCount = 1
			}
			if pageCount >= sa.startPage && pageCount <= sa.endPage {
				fmt.Fprintf(fileout, "%s", line)
			}
		}
	} else {
		for {
			page, err := readLine.ReadString('\f')
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: Read file error!\n", progname)
				os.Exit(8)
			}
			pageCount++
			if pageCount >= sa.startPage && pageCount <= sa.endPage {
				fmt.Fprintf(fileout, "%s", page)
			}
		}
	}
	cmd := exec.Command("cat", "-n")
	_, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: Create pipe error\n", progname)
		os.Exit(9)
	}
	if sa.printDest != "" {
		cmd.Stdout = fileout
		cmd.Run()
	}
	filein.Close()
	fileout.Close()
}

func main() {
	sa := selpgArgs{0, 0, "", 7, false, ""}
	progname = os.Args[0]
	ProcessArgs(&sa)
	ProcessInput(&sa)
}
