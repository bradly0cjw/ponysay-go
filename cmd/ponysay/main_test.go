package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCLIFullFeatureSuite(t *testing.T) {
	// 1. Build test binary
	buildCmd := exec.Command("go", "build", "-o", "ponysay_test_bin", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build test binary: %v", err)
	}
	defer exec.Command("rm", "-f", "ponysay_test_bin", "ponythink_test_bin").Run()

	// Symlink for ponythink
	exec.Command("ln", "-sf", "ponysay_test_bin", "ponythink_test_bin").Run()

	// 2. Test Basic message execution
	cmd := exec.Command("./ponysay_test_bin", "Hello Test!")
	out, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(out), "Hello Test!") {
		t.Errorf("Basic message failed. Output: %s", string(out))
	}

	// 3. Test Pony selection (-f)
	fCmd := exec.Command("./ponysay_test_bin", "-f", "pinkie", "Party!")
	fOut, fErr := fCmd.CombinedOutput()
	if fErr != nil || len(fOut) == 0 {
		t.Errorf("Pony selection -f failed")
	}

	// 4. Test Non-MLP Pony (+f)
	plusFCmd := exec.Command("./ponysay_test_bin", "+f", "cow", "Moo!")
	plusFOut, _ := plusFCmd.CombinedOutput()
	if len(plusFOut) == 0 {
		t.Errorf("Non-MLP pony +f failed")
	}

	// 5. Test Quote mode (-q pinkie)
	qCmd := exec.Command("./ponysay_test_bin", "-q", "pinkie")
	qOut, qErr := qCmd.CombinedOutput()
	if qErr != nil || len(qOut) == 0 {
		t.Errorf("Quote mode -q pinkie failed")
	}

	// 6. Test Quoters listing (--quoters)
	quotersCmd := exec.Command("./ponysay_test_bin", "--quoters")
	qList, qListErr := quotersCmd.CombinedOutput()
	if qListErr != nil || !strings.Contains(string(qList), "pinkie") {
		t.Errorf("--quoters failed")
	}

	// 7. Test Pony listing (-l, -L, +l, +L, -A, +A)
	lCmd := exec.Command("./ponysay_test_bin", "-l")
	lOut, _ := lCmd.CombinedOutput()
	if !strings.Contains(string(lOut), "derpy") {
		t.Errorf("-l listing failed")
	}

	lAliasesCmd := exec.Command("./ponysay_test_bin", "-L")
	lAliasesOut, _ := lAliasesCmd.CombinedOutput()
	if !strings.Contains(string(lAliasesOut), "derpy") {
		t.Errorf("-L listing with aliases failed")
	}

	// 8. Test One-line listing (--onelist)
	oneListCmd := exec.Command("./ponysay_test_bin", "-l", "--onelist")
	oneListOut, _ := oneListCmd.CombinedOutput()
	if strings.Count(string(oneListOut), "\n") > 2 {
		t.Errorf("--onelist should output list in a single line")
	}

	// 9. Test Balloon listing (-B)
	bListCmd := exec.Command("./ponysay_test_bin", "-B")
	bListOut, _ := bListCmd.CombinedOutput()
	if !strings.Contains(string(bListOut), "unicode") || !strings.Contains(string(bListOut), "cowsay") {
		t.Errorf("-B balloon listing failed")
	}

	// 10. Test Metadata Info (-i and +i)
	infoCmd := exec.Command("./ponysay_test_bin", "-f", "derpy", "-i")
	infoOut, _ := infoCmd.CombinedOutput()
	if !strings.Contains(string(infoOut), "NAME:") {
		t.Errorf("-i metadata info failed")
	}

	// 11. Test Pony Only (-o)
	ponyOnlyCmd := exec.Command("./ponysay_test_bin", "-f", "derpy", "-o")
	ponyOnlyOut, _ := ponyOnlyCmd.CombinedOutput()
	if strings.Contains(string(ponyOnlyOut), "┌") || strings.Contains(string(ponyOnlyOut), "<") {
		t.Errorf("-o pony-only failed, balloon was rendered")
	}

	// 12. Test Ponythink mode
	thinkCmd := exec.Command("./ponythink_test_bin", "Deep thoughts...")
	thinkOut, _ := thinkCmd.CombinedOutput()
	if !strings.Contains(string(thinkOut), "(") && !strings.Contains(string(thinkOut), "o") {
		t.Errorf("ponythink failed")
	}

	// 13. Test Help output containing update command
	helpCmd := exec.Command("./ponysay_test_bin", "-h")
	helpOut, _ := helpCmd.CombinedOutput()
	if !strings.Contains(string(helpOut), "update") {
		t.Errorf("-h help output missing update command description")
	}

	// 14. Test -F (any pony) flag
	anyCmd := exec.Command("./ponysay_test_bin", "-F", "fluttershy", "AnyPony Test")
	anyOut, anyErr := anyCmd.CombinedOutput()
	if anyErr != nil || !strings.Contains(string(anyOut), "AnyPony Test") {
		t.Errorf("-F flag failed: %v, output: %s", anyErr, string(anyOut))
	}

	// 15. Test Fuzzy Matching for misspelled pony
	fuzzyCmd := exec.Command("./ponysay_test_bin", "-f", "fluter-shy", "Fuzzy Test")
	fuzzyOut, fuzzyErr := fuzzyCmd.CombinedOutput()
	if fuzzyErr != nil || !strings.Contains(string(fuzzyOut), "Fuzzy Test") {
		t.Errorf("Fuzzy pony search failed: %v, output: %s", fuzzyErr, string(fuzzyOut))
	}

	// 16. Test Attached flag syntax (-ffluttershy)
	attCmd := exec.Command("./ponysay_test_bin", "-ffluttershy", "Attached Test")
	attOut, attErr := attCmd.CombinedOutput()
	if attErr != nil || !strings.Contains(string(attOut), "Attached Test") {
		t.Errorf("Attached flag -ffluttershy failed: %v, output: %s", attErr, string(attOut))
	}

	// 17. Test argument placement flexibility (message before flags)
	posCmd := exec.Command("./ponysay_test_bin", "Positional Test", "-f", "fluttershy", "-b", "cowsay")
	posOut, posErr := posCmd.CombinedOutput()
	if posErr != nil || !strings.Contains(string(posOut), "Positional Test") {
		t.Errorf("Positional argument before flags failed: %v, output: %s", posErr, string(posOut))
	}

	// 18. Test Fuzzy Quote Selection (-q fluter-shy)
	fqCmd := exec.Command("./ponysay_test_bin", "-q", "fluter-shy")
	fqOut, fqErr := fqCmd.CombinedOutput()
	if fqErr != nil || strings.Contains(string(fqOut), "mute!") || len(fqOut) == 0 {
		t.Errorf("Fuzzy quote selection failed: %v, output: %s", fqErr, string(fqOut))
	}
}

