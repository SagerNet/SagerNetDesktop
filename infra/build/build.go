package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type buildType struct {
	goos   string
	goarch string
	goarm  string
}

func main() {
	_, gnf := exec.LookPath("go")
	if gnf != nil {
		path, _ := os.LookupEnv("PATH")
		_ = os.Setenv("PATH", fmt.Sprint(path, ":/usr/lib/go-1.16/bin"))
	}

	version, _ := os.LookupEnv("VERSION")
	if version == "" {
		panic("missing version")
	}

	var builds []buildType
	isAll := false

	if len(os.Args) > 1 && os.Args[1] == "all" {

		isAll = true
		goos := []string{"windows" /*"freebsd", "openbsd",*/, "linux", "android"}
		goarch := []string{"amd64", "386"}
		for _, o := range goos {
			for _, a := range goarch {
				builds = append(builds, buildType{goos: o, goarch: a})
			}
		}

		builds = append(builds, buildType{goos: "windows", goarch: "arm", goarm: "5"})
		builds = append(builds, buildType{goos: "windows", goarch: "arm", goarm: "6"})
		builds = append(builds, buildType{goos: "windows", goarch: "arm", goarm: "7"})

		builds = append(builds, buildType{goos: "linux", goarch: "arm", goarm: "5"})
		builds = append(builds, buildType{goos: "linux", goarch: "arm", goarm: "6"})
		builds = append(builds, buildType{goos: "linux", goarch: "arm", goarm: "7"})

		builds = append(builds, buildType{goos: "linux", goarch: "arm64"})
		//builds = append(builds, buildType{goos: "linux", goarch: "riscv64"})

		builds = append(builds, buildType{goos: "linux", goarch: "mips64"})
		builds = append(builds, buildType{goos: "linux", goarch: "mips64le"})
		builds = append(builds, buildType{goos: "linux", goarch: "mips"})

		builds = append(builds, buildType{goos: "android", goarch: "arm", goarm: "5"})
		builds = append(builds, buildType{goos: "android", goarch: "arm", goarm: "6"})
		builds = append(builds, buildType{goos: "android", goarch: "arm", goarm: "7"})
		builds = append(builds, buildType{goos: "android", goarch: "arm64"})

		builds = append(builds, buildType{goos: "darwin", goarch: "amd64"})
		builds = append(builds, buildType{goos: "darwin", goarch: "arm64"})
	} else {
		goos, _ := os.LookupEnv("GOOS")
		goarch, _ := os.LookupEnv("GOARCH")
		goarm, _ := os.LookupEnv("GOARM")

		if goos == "" {
			goos = runtime.GOOS
		}

		if goarch == "" {
			goarch = runtime.GOARCH
		}

		builds = append(builds, buildType{goos: goos, goarch: goarch, goarm: goarm})
	}

	_, buildAndroid := os.LookupEnv("ANDROID_ARM_CC")

	for _, build := range builds {
		arch := build.goarch
		switch arch {
		case "arm":
			arch = fmt.Sprint(arch, "-", build.goarm)
			break
		case "arm64":
			arch = "aarch64"
			break
		case "386":
			arch = "x86"
			break
		case "amd64":
			arch = "x86_64"
			break
		}

		if build.goos == "android" {
			if !buildAndroid {
				continue
			}
		} else if build.goos == "darwin" {
			xgo, _ := os.LookupEnv("XGO")
			if xgo == "1" {
				output := fmt.Sprint("build/sagernet-", version, "-", build.goos, "-", arch)
				shell := exec.Command("xgo", "-v", "-go", "go-1.16.7", "-dest", output, "-targets", fmt.Sprint(build.goos, "/", build.goarch), "-out", "sagernet", ".")
				shell.Stdin = os.Stdin
				shell.Stdout = os.Stdout
				shell.Stderr = os.Stderr

				err := shell.Run()
				if err != nil {
					println(err.Error())
					os.Exit(shell.ProcessState.ExitCode())
				}

				err = exec.Command("sudo", "chmod", "-R", "777", output).Run()
				if err != nil {
					panic(err.Error())
				}

				dir, err := os.Open(output)
				if err != nil {
					panic(err)
				}

				binary, err := dir.ReadDir(1)
				if err != nil {
					panic(err)
				}

				_ = dir.Close()
				err = os.Rename(fmt.Sprint(output, "/", binary[0].Name()), fmt.Sprint(output, "/sagernet"))
				if err != nil {
					panic(err)
				}
				continue
			}
		}

		var output string
		if isAll {
			output = fmt.Sprint("build/sagernet-", version, "-", build.goos, "-", arch, "/sagernet")
		} else {
			output = "build/current/sagernet"
		}

		if build.goos == "windows" {
			output = fmt.Sprint(output, ".exe")
		}

		args := []string{"build", "-o", output, "-v", "-trimpath"}

		_ = os.Setenv("GOOS", build.goos)
		_ = os.Setenv("GOARCH", build.goarch)
		if build.goarch == "arm" {
			_ = os.Setenv("GOARM", build.goarm)
		} else {
			_ = os.Unsetenv("GOARM")
		}

		println(">", build.goos, arch)

		if build.goos == "windows" {
			//args = append(args, "-ldflags -H=windowsgui -s -w -buildid=")

			switch build.goarch {
			case "amd64":
				_ = os.Setenv("CC", "x86_64-w64-mingw32-gcc")
				break
			case "386":
				_ = os.Setenv("CC", "i686-w64-mingw32-gcc")
				break
			}
		} else if build.goos == "linux" {
			switch build.goarch {
			case "amd64":
				_ = os.Setenv("CC", "x86_64-linux-gnu-gcc")
				break
			case "386":
				_ = os.Setenv("CC", "i686-linux-gnu-gcc")
				break
			case "arm":
				_ = os.Setenv("CC", "arm-linux-gnueabi-gcc")
				break
			case "arm64":
				_ = os.Setenv("CC", "aarch64-linux-gnu-gcc")
				break
			case "mips":
				_ = os.Setenv("CC", "mips-linux-gnu-gcc")
				break
			case "mips64":
				_ = os.Setenv("CC", "mips64-linux-gnuabi64-gcc")
				break
			case "mips64le":
				_ = os.Setenv("CC", "mips64el-linux-gnuabi64-gcc")
				break
			case "risc64":
				_ = os.Setenv("CC", "riscv64-linux-gnu-gcc")
				break
			}
		} else if build.goos == "android" {
			switch build.goarch {
			case "arm":
				cc, _ := os.LookupEnv("ANDROID_ARM_CC")
				_ = os.Setenv("CC", cc)
				break
			case "arm64":
				cc, _ := os.LookupEnv("ANDROID_ARM64_CC")
				_ = os.Setenv("CC", cc)
				break
			case "amd64":
				cc, _ := os.LookupEnv("ANDROID_X86_64_CC")
				_ = os.Setenv("CC", cc)
				break
			case "386":
				cc, _ := os.LookupEnv("ANDROID_X86_CC")
				_ = os.Setenv("CC", cc)
				break
			}
		} else {
			_ = os.Unsetenv("CC")
		}
		args = append(args, "-ldflags", "-s -w -buildid=")
		args = append(args, ".")

		println(">", fmt.Sprint("go ", strings.Join(args, " ")))

		_ = os.Setenv("CGO_ENABLED", "1")

		shell := exec.Command("go", args...)
		shell.Stdin = os.Stdin
		shell.Stdout = os.Stdout
		shell.Stderr = os.Stderr

		err := shell.Run()
		if err != nil {
			println(err.Error())
			os.Exit(shell.ProcessState.ExitCode())
		}
	}

	if isAll {

		// make zip or tar

		os.RemoveAll("build/release")
		err := os.Mkdir("build/release", 0777)
		if err != nil {
			panic(err.Error())
		}

		outDir, _ := os.Open("build")
		dirEntries, _ := outDir.ReadDir(114)
		for _, subDir := range dirEntries {
			if !subDir.IsDir() || subDir.Name() == "currnet" || subDir.Name() == "release" {
				continue
			}

			if strings.Contains(subDir.Name(), "windows") {
				// use zip
				inFile, err := os.Open(fmt.Sprint("build/", subDir.Name(), "/sagernet.exe"))
				if err != nil {
					panic(err.Error())
				}

				var outFile *os.File
				outPath := fmt.Sprint("build/release/", subDir.Name(), ".zip")

				_ = os.RemoveAll(outPath)
				outFile, err = os.Create(outPath)
				if err != nil {
					panic(err.Error())
				}

				writer := zip.NewWriter(outFile)
				output, err := writer.Create("sagernet.exe")
				if err != nil {
					panic(err.Error())
				}

				_, err = io.Copy(output, inFile)
				if err != nil {
					panic(err.Error())
				}

				_ = writer.Flush()
				_ = writer.Close()
				_ = outFile.Close()
				_ = inFile.Close()
			} else {
				// use tar.gz

				inFile, err := os.Open(fmt.Sprint("build/", subDir.Name(), "/sagernet"))
				if err != nil {
					panic(err.Error())
				}

				inInfo, err := inFile.Stat()
				if err != nil {
					panic(err.Error())
				}

				var outFile *os.File

				outPath := fmt.Sprint("build/release/", subDir.Name(), ".tar.gz")
				_ = os.RemoveAll(outPath)

				outFile, err = os.Create(outPath)
				gzWriter := gzip.NewWriter(outFile)
				writer := tar.NewWriter(gzWriter)
				header, err := tar.FileInfoHeader(inInfo, "")
				err = writer.WriteHeader(header)
				if err != nil {
					panic(err.Error())
				}

				_, err = io.Copy(writer, inFile)
				if err != nil {
					panic(err.Error())
				}

				_ = writer.Flush()
				_ = writer.Close()
				_ = gzWriter.Flush()
				_ = gzWriter.Close()
				_ = outFile.Close()
				_ = inFile.Close()
			}
		}

		_ = outDir.Close()

	}

}
