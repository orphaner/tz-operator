package k8s

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"testing"
)

var _ = Describe("ToDNSName", func() {
	table.DescribeTable("is DNS name compliant",
		func(prefix string, name string, suffix string, expected string) {
			result := ToDNSName(prefix, name, suffix)
			Expect(result).To(Equal(expected))
			Expect(len(result) <= 63).Should(BeTrue())
		},
		Entry("<= 63", "kidle", "shortname", "idle", "kidle-shortname-idle"),
		Entry("<= 63", "kidle", "shortname", "", "kidle-shortname"),
		Entry("== 63", "kidle", "name-length-is-63-yessssssssssssssssssssssssssssss", "wakeup", "kidle-name-length-is-63-yessssssssssssssssssssssssssssss-wakeup"),
		Entry("> 63", "kidle", "very-toooooooooooooooooooooooooooooooooooooooooooooooooo-long", "idle", "kidle-very-tooooooooooooooooooooooooooooooooooooooo-dmvyes-idle"),
	)
})

func TestBase64(t *testing.T) {
	RegisterTestingT(t)

	cases := []struct {
		name           string
		s              string
		length         int
		expected       string
		expectedLength int
	}{
		{"exact length", "kidle is awesome", 24, "a2lkbGUgaXMgYXdlc29tZQ==", 24},
		{"over lengthed", "kidle is awesome", 666, "a2lkbGUgaXMgYXdlc29tZQ==", 24},
		{"length == 1", "kidle is awesome", 1, "a", 1},
		{"length == 0", "kidle is awesome", 0, "", 0},
	}
	for _, tc := range cases {
		b := Base64(tc.s, tc.length)
		Expect(b).To(Equal(tc.expected))
		Expect(len(b)).Should(BeIdenticalTo(tc.expectedLength))
	}
}
