package lemin

import (
	"crypto/sha256"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// Digests captured from the original CLI at 630309b, before parser extraction.
// Include the map echo, blank lines, move ordering and trailing spaces.
func TestOriginalCLIStdout(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(t.TempDir(), "lem-in")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	cases := map[string]string{
		"input.txt":          "c31941a7a673626fdeac2550fef50bd87d2d4c27189b620ebd107ad78141ee61",
		"3D/input.txt":       "6936e2cfd9906f69c93e2e733afbd8efe5e5131fca67b2c1f0ea6de4ad88e793",
		"test/example00.txt": "c31941a7a673626fdeac2550fef50bd87d2d4c27189b620ebd107ad78141ee61",
		"test/example01.txt": "41d70ddc30a7d43827be6b197ee64bf3dd1f160be17db1b772c76bb015fa097e",
		"test/example02.txt": "761ab05ec06092de6d15551e72844a229dc77ffa2f2593e8dececb1c78a0016b",
		"test/example03.txt": "d7e60d8049e36df9f46fde8b4adb590b37f90a620b0a465f5d9ccdca4c3c3da8",
		"test/example04.txt": "508a9312752b11a2e73bf6265767a1c360dfbab90af9cc9ac86fea3a749a8c68",
		"test/example05.txt": "62c999ee6415948275a2637a9c693e3f18ad8ba17597d56f98a470afe59ed566",
		"test/example06.txt": "d554982357ecc8fcb29c44f1a91427165f724ad52923f3b002125fa468625344",
		"test/example07.txt": "c1daa21e4a6e3c64a091c79ee7410c53c4cda24ba5bb0369347b3b55e5fb250d",
	}
	for file, want := range cases {
		t.Run(file, func(t *testing.T) {
			command := exec.Command(bin, filepath.Join(root, filepath.FromSlash(file)))
			command.Dir = t.TempDir()
			out, err := command.Output()
			if err != nil {
				t.Fatal(err)
			}
			if got := fmt.Sprintf("%x", sha256.Sum256(out)); got != want {
				t.Fatalf("original CLI stdout changed: sha256=%s, want %s", got, want)
			}
		})
	}
	command := exec.Command(bin, filepath.Join(root, "test", "badexample00.txt"))
	command.Dir = t.TempDir()
	out, err := command.Output()
	if err == nil || string(out) != "ERROR: invalid data format\n" {
		t.Fatalf("invalid map rejection changed: err=%v stdout=%q", err, out)
	}
}
