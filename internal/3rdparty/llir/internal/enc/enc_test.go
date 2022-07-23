package enc

import (
	"reflect"
	"testing"
)

func TestGlobalName(t *testing.T) {
	golden := []struct {
		s    string
		want string
	}{
		// i=0
		{s: "foo", want: "@foo"},
		// i=1
		{s: "a b", want: `@"a b"`},
		// i=2
		{s: "$a", want: "@$a"},
		// i=3
		{s: "-a", want: "@-a"},
		// i=4
		{s: ".a", want: "@.a"},
		// i=5
		{s: "_a", want: "@_a"},
		// i=6
		{s: "#a", want: `@"#a"`},
		// i=7
		{s: "a b#c", want: `@"a b#c"`},
		// i=8
		{s: "2", want: `@"2"`},
		// i=9
		{s: "foo世bar", want: `@"foo\E4\B8\96bar"`},
		// i=10
		{s: "foo\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29\x2A\x2B\x2C\x2D\x2E\x2F\x30\x31\x32\x33\x34\x35\x36\x37\x38\x39\x3A\x3B\x3C\x3D\x3E\x3F\x40\x41\x42\x43\x44\x45\x46\x47\x48\x49\x4A\x4B\x4C\x4D\x4E\x4F\x50\x51\x52\x53\x54\x55\x56\x57\x58\x59\x5A\x5B\x5C\x5D\x5E\x5F\x60\x61\x62\x63\x64\x65\x66\x67\x68\x69\x6A\x6B\x6C\x6D\x6E\x6F\x70\x71\x72\x73\x74\x75\x76\x77\x78\x79\x7A\x7B\x7C\x7D\x7E\x7F\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8A\x8B\x8C\x8D\x8E\x8F\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9A\x9B\x9C\x9D\x9E\x9F\xA0\xA1\xA2\xA3\xA4\xA5\xA6\xA7\xA8\xA9\xAA\xAB\xAC\xAD\xAE\xAF\xB0\xB1\xB2\xB3\xB4\xB5\xB6\xB7\xB8\xB9\xBA\xBB\xBC\xBD\xBE\xBF\xC0\xC1\xC2\xC3\xC4\xC5\xC6\xC7\xC8\xC9\xCA\xCB\xCC\xCD\xCE\xCF\xD0\xD1\xD2\xD3\xD4\xD5\xD6\xD7\xD8\xD9\xDA\xDB\xDC\xDD\xDE\xDF\xE0\xE1\xE2\xE3\xE4\xE5\xE6\xE7\xE8\xE9\xEA\xEB\xEC\xED\xEE\xEF\xF0\xF1\xF2\xF3\xF4\xF5\xF6\xF7\xF8\xF9\xFA\xFB\xFC\xFD\xFE\xFF", want: "@\"foo\\01\\02\\03\\04\\05\\06\\07\\08\\09\\0A\\0B\\0C\\0D\\0E\\0F\\10\\11\\12\\13\\14\\15\\16\\17\\18\\19\\1A\\1B\\1C\\1D\\1E\\1F !\\22#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\5C]^_`abcdefghijklmnopqrstuvwxyz{|}~\\7F\\80\\81\\82\\83\\84\\85\\86\\87\\88\\89\\8A\\8B\\8C\\8D\\8E\\8F\\90\\91\\92\\93\\94\\95\\96\\97\\98\\99\\9A\\9B\\9C\\9D\\9E\\9F\\A0\\A1\\A2\\A3\\A4\\A5\\A6\\A7\\A8\\A9\\AA\\AB\\AC\\AD\\AE\\AF\\B0\\B1\\B2\\B3\\B4\\B5\\B6\\B7\\B8\\B9\\BA\\BB\\BC\\BD\\BE\\BF\\C0\\C1\\C2\\C3\\C4\\C5\\C6\\C7\\C8\\C9\\CA\\CB\\CC\\CD\\CE\\CF\\D0\\D1\\D2\\D3\\D4\\D5\\D6\\D7\\D8\\D9\\DA\\DB\\DC\\DD\\DE\\DF\\E0\\E1\\E2\\E3\\E4\\E5\\E6\\E7\\E8\\E9\\EA\\EB\\EC\\ED\\EE\\EF\\F0\\F1\\F2\\F3\\F4\\F5\\F6\\F7\\F8\\F9\\FA\\FB\\FC\\FD\\FE\\FF\""},
	}
	for i, g := range golden {
		got := GlobalName(g.s)
		if g.want != got {
			t.Errorf("i=%d: name mismatch; expected `%s`, got `%s`", i, g.want, got)
		}
	}
}

func TestLocalName(t *testing.T) {
	golden := []struct {
		s    string
		want string
	}{
		// i=0
		{s: "foo", want: "%foo"},
		// i=1
		{s: "a b", want: `%"a b"`},
		// i=2
		{s: "$a", want: "%$a"},
		// i=3
		{s: "-a", want: "%-a"},
		// i=4
		{s: ".a", want: "%.a"},
		// i=5
		{s: "_a", want: "%_a"},
		// i=6
		{s: "#a", want: `%"#a"`},
		// i=7
		{s: "a b#c", want: `%"a b#c"`},
		// i=8
		{s: "2", want: `%"2"`},
		// i=9
		{s: "foo世bar", want: `%"foo\E4\B8\96bar"`},
		// i=10
		{s: "foo\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29\x2A\x2B\x2C\x2D\x2E\x2F\x30\x31\x32\x33\x34\x35\x36\x37\x38\x39\x3A\x3B\x3C\x3D\x3E\x3F\x40\x41\x42\x43\x44\x45\x46\x47\x48\x49\x4A\x4B\x4C\x4D\x4E\x4F\x50\x51\x52\x53\x54\x55\x56\x57\x58\x59\x5A\x5B\x5C\x5D\x5E\x5F\x60\x61\x62\x63\x64\x65\x66\x67\x68\x69\x6A\x6B\x6C\x6D\x6E\x6F\x70\x71\x72\x73\x74\x75\x76\x77\x78\x79\x7A\x7B\x7C\x7D\x7E\x7F\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8A\x8B\x8C\x8D\x8E\x8F\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9A\x9B\x9C\x9D\x9E\x9F\xA0\xA1\xA2\xA3\xA4\xA5\xA6\xA7\xA8\xA9\xAA\xAB\xAC\xAD\xAE\xAF\xB0\xB1\xB2\xB3\xB4\xB5\xB6\xB7\xB8\xB9\xBA\xBB\xBC\xBD\xBE\xBF\xC0\xC1\xC2\xC3\xC4\xC5\xC6\xC7\xC8\xC9\xCA\xCB\xCC\xCD\xCE\xCF\xD0\xD1\xD2\xD3\xD4\xD5\xD6\xD7\xD8\xD9\xDA\xDB\xDC\xDD\xDE\xDF\xE0\xE1\xE2\xE3\xE4\xE5\xE6\xE7\xE8\xE9\xEA\xEB\xEC\xED\xEE\xEF\xF0\xF1\xF2\xF3\xF4\xF5\xF6\xF7\xF8\xF9\xFA\xFB\xFC\xFD\xFE\xFF", want: "%\"foo\\01\\02\\03\\04\\05\\06\\07\\08\\09\\0A\\0B\\0C\\0D\\0E\\0F\\10\\11\\12\\13\\14\\15\\16\\17\\18\\19\\1A\\1B\\1C\\1D\\1E\\1F !\\22#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\5C]^_`abcdefghijklmnopqrstuvwxyz{|}~\\7F\\80\\81\\82\\83\\84\\85\\86\\87\\88\\89\\8A\\8B\\8C\\8D\\8E\\8F\\90\\91\\92\\93\\94\\95\\96\\97\\98\\99\\9A\\9B\\9C\\9D\\9E\\9F\\A0\\A1\\A2\\A3\\A4\\A5\\A6\\A7\\A8\\A9\\AA\\AB\\AC\\AD\\AE\\AF\\B0\\B1\\B2\\B3\\B4\\B5\\B6\\B7\\B8\\B9\\BA\\BB\\BC\\BD\\BE\\BF\\C0\\C1\\C2\\C3\\C4\\C5\\C6\\C7\\C8\\C9\\CA\\CB\\CC\\CD\\CE\\CF\\D0\\D1\\D2\\D3\\D4\\D5\\D6\\D7\\D8\\D9\\DA\\DB\\DC\\DD\\DE\\DF\\E0\\E1\\E2\\E3\\E4\\E5\\E6\\E7\\E8\\E9\\EA\\EB\\EC\\ED\\EE\\EF\\F0\\F1\\F2\\F3\\F4\\F5\\F6\\F7\\F8\\F9\\FA\\FB\\FC\\FD\\FE\\FF\""},
	}
	for i, g := range golden {
		got := LocalName(g.s)
		if g.want != got {
			t.Errorf("i=%d: name mismatch; expected `%s`, got `%s`", i, g.want, got)
		}
	}
}

func TestLabelName(t *testing.T) {
	golden := []struct {
		s    string
		want string
	}{
		// i=0
		{s: "foo", want: "foo:"},
		// i=1
		{s: "a b", want: `"a b":`},
		// i=2
		{s: "$a", want: "$a:"},
		// i=3
		{s: "-a", want: "-a:"},
		// i=4
		{s: ".a", want: ".a:"},
		// i=5
		{s: "_a", want: "_a:"},
		// i=6
		{s: "#a", want: `"#a":`},
		// i=7
		{s: "a b#c", want: `"a b#c":`},
		// i=8
		{s: "2", want: `"2":`},
		// i=9
		{s: "foo世bar", want: `"foo\E4\B8\96bar":`},
		// i=10
		{s: "foo\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29\x2A\x2B\x2C\x2D\x2E\x2F\x30\x31\x32\x33\x34\x35\x36\x37\x38\x39\x3A\x3B\x3C\x3D\x3E\x3F\x40\x41\x42\x43\x44\x45\x46\x47\x48\x49\x4A\x4B\x4C\x4D\x4E\x4F\x50\x51\x52\x53\x54\x55\x56\x57\x58\x59\x5A\x5B\x5C\x5D\x5E\x5F\x60\x61\x62\x63\x64\x65\x66\x67\x68\x69\x6A\x6B\x6C\x6D\x6E\x6F\x70\x71\x72\x73\x74\x75\x76\x77\x78\x79\x7A\x7B\x7C\x7D\x7E\x7F\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8A\x8B\x8C\x8D\x8E\x8F\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9A\x9B\x9C\x9D\x9E\x9F\xA0\xA1\xA2\xA3\xA4\xA5\xA6\xA7\xA8\xA9\xAA\xAB\xAC\xAD\xAE\xAF\xB0\xB1\xB2\xB3\xB4\xB5\xB6\xB7\xB8\xB9\xBA\xBB\xBC\xBD\xBE\xBF\xC0\xC1\xC2\xC3\xC4\xC5\xC6\xC7\xC8\xC9\xCA\xCB\xCC\xCD\xCE\xCF\xD0\xD1\xD2\xD3\xD4\xD5\xD6\xD7\xD8\xD9\xDA\xDB\xDC\xDD\xDE\xDF\xE0\xE1\xE2\xE3\xE4\xE5\xE6\xE7\xE8\xE9\xEA\xEB\xEC\xED\xEE\xEF\xF0\xF1\xF2\xF3\xF4\xF5\xF6\xF7\xF8\xF9\xFA\xFB\xFC\xFD\xFE\xFF", want: "\"foo\\01\\02\\03\\04\\05\\06\\07\\08\\09\\0A\\0B\\0C\\0D\\0E\\0F\\10\\11\\12\\13\\14\\15\\16\\17\\18\\19\\1A\\1B\\1C\\1D\\1E\\1F !\\22#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\5C]^_`abcdefghijklmnopqrstuvwxyz{|}~\\7F\\80\\81\\82\\83\\84\\85\\86\\87\\88\\89\\8A\\8B\\8C\\8D\\8E\\8F\\90\\91\\92\\93\\94\\95\\96\\97\\98\\99\\9A\\9B\\9C\\9D\\9E\\9F\\A0\\A1\\A2\\A3\\A4\\A5\\A6\\A7\\A8\\A9\\AA\\AB\\AC\\AD\\AE\\AF\\B0\\B1\\B2\\B3\\B4\\B5\\B6\\B7\\B8\\B9\\BA\\BB\\BC\\BD\\BE\\BF\\C0\\C1\\C2\\C3\\C4\\C5\\C6\\C7\\C8\\C9\\CA\\CB\\CC\\CD\\CE\\CF\\D0\\D1\\D2\\D3\\D4\\D5\\D6\\D7\\D8\\D9\\DA\\DB\\DC\\DD\\DE\\DF\\E0\\E1\\E2\\E3\\E4\\E5\\E6\\E7\\E8\\E9\\EA\\EB\\EC\\ED\\EE\\EF\\F0\\F1\\F2\\F3\\F4\\F5\\F6\\F7\\F8\\F9\\FA\\FB\\FC\\FD\\FE\\FF\":"},
	}
	for i, g := range golden {
		got := LabelName(g.s)
		if g.want != got {
			t.Errorf("i=%d: name mismatch; expected `%s`, got `%s`", i, g.want, got)
		}
	}
}

func TestTypeName(t *testing.T) {
	golden := []struct {
		s    string
		want string
	}{
		// i=0
		{s: "foo", want: "%foo"},
		// i=1
		{s: "a b", want: `%"a b"`},
		// i=2
		{s: "$a", want: "%$a"},
		// i=3
		{s: "-a", want: "%-a"},
		// i=4
		{s: ".a", want: "%.a"},
		// i=5
		{s: "_a", want: "%_a"},
		// i=6
		{s: "#a", want: `%"#a"`},
		// i=7
		{s: "a b#c", want: `%"a b#c"`},
		// i=8
		{s: "2", want: `%2`}, // the only difference between enc.LocalName and enc.TypeName
		// i=9
		{s: "foo世bar", want: `%"foo\E4\B8\96bar"`},
		// i=10
		{s: "foo\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29\x2A\x2B\x2C\x2D\x2E\x2F\x30\x31\x32\x33\x34\x35\x36\x37\x38\x39\x3A\x3B\x3C\x3D\x3E\x3F\x40\x41\x42\x43\x44\x45\x46\x47\x48\x49\x4A\x4B\x4C\x4D\x4E\x4F\x50\x51\x52\x53\x54\x55\x56\x57\x58\x59\x5A\x5B\x5C\x5D\x5E\x5F\x60\x61\x62\x63\x64\x65\x66\x67\x68\x69\x6A\x6B\x6C\x6D\x6E\x6F\x70\x71\x72\x73\x74\x75\x76\x77\x78\x79\x7A\x7B\x7C\x7D\x7E\x7F\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8A\x8B\x8C\x8D\x8E\x8F\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9A\x9B\x9C\x9D\x9E\x9F\xA0\xA1\xA2\xA3\xA4\xA5\xA6\xA7\xA8\xA9\xAA\xAB\xAC\xAD\xAE\xAF\xB0\xB1\xB2\xB3\xB4\xB5\xB6\xB7\xB8\xB9\xBA\xBB\xBC\xBD\xBE\xBF\xC0\xC1\xC2\xC3\xC4\xC5\xC6\xC7\xC8\xC9\xCA\xCB\xCC\xCD\xCE\xCF\xD0\xD1\xD2\xD3\xD4\xD5\xD6\xD7\xD8\xD9\xDA\xDB\xDC\xDD\xDE\xDF\xE0\xE1\xE2\xE3\xE4\xE5\xE6\xE7\xE8\xE9\xEA\xEB\xEC\xED\xEE\xEF\xF0\xF1\xF2\xF3\xF4\xF5\xF6\xF7\xF8\xF9\xFA\xFB\xFC\xFD\xFE\xFF", want: "%\"foo\\01\\02\\03\\04\\05\\06\\07\\08\\09\\0A\\0B\\0C\\0D\\0E\\0F\\10\\11\\12\\13\\14\\15\\16\\17\\18\\19\\1A\\1B\\1C\\1D\\1E\\1F !\\22#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\5C]^_`abcdefghijklmnopqrstuvwxyz{|}~\\7F\\80\\81\\82\\83\\84\\85\\86\\87\\88\\89\\8A\\8B\\8C\\8D\\8E\\8F\\90\\91\\92\\93\\94\\95\\96\\97\\98\\99\\9A\\9B\\9C\\9D\\9E\\9F\\A0\\A1\\A2\\A3\\A4\\A5\\A6\\A7\\A8\\A9\\AA\\AB\\AC\\AD\\AE\\AF\\B0\\B1\\B2\\B3\\B4\\B5\\B6\\B7\\B8\\B9\\BA\\BB\\BC\\BD\\BE\\BF\\C0\\C1\\C2\\C3\\C4\\C5\\C6\\C7\\C8\\C9\\CA\\CB\\CC\\CD\\CE\\CF\\D0\\D1\\D2\\D3\\D4\\D5\\D6\\D7\\D8\\D9\\DA\\DB\\DC\\DD\\DE\\DF\\E0\\E1\\E2\\E3\\E4\\E5\\E6\\E7\\E8\\E9\\EA\\EB\\EC\\ED\\EE\\EF\\F0\\F1\\F2\\F3\\F4\\F5\\F6\\F7\\F8\\F9\\FA\\FB\\FC\\FD\\FE\\FF\""},
	}
	for i, g := range golden {
		got := TypeName(g.s)
		if g.want != got {
			t.Errorf("i=%d: name mismatch; expected `%s`, got `%s`", i, g.want, got)
		}
	}
}

func TestAttrGroupID(t *testing.T) {
	golden := []struct {
		s    int64
		want string
	}{
		// i=0
		{s: 42, want: "#42"},
	}
	for i, g := range golden {
		got := AttrGroupID(g.s)
		if g.want != got {
			t.Errorf("i=%d: name mismatch; expected `%s`, got `%s`", i, g.want, got)
		}
	}
}

func TestComdatName(t *testing.T) {
	golden := []struct {
		s    string
		want string
	}{
		// i=0
		{s: "foo", want: "$foo"},
		// i=1
		{s: "a b", want: `$"a b"`},
		// i=2
		{s: "$a", want: "$$a"},
		// i=3
		{s: "-a", want: "$-a"},
		// i=4
		{s: ".a", want: "$.a"},
		// i=5
		{s: "_a", want: "$_a"},
		// i=6
		{s: "#a", want: `$"#a"`},
		// i=7
		{s: "a b#c", want: `$"a b#c"`},
		// i=8
		{s: "foo世bar", want: `$"foo\E4\B8\96bar"`},
		// i=9
		{s: "foo\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29\x2A\x2B\x2C\x2D\x2E\x2F\x30\x31\x32\x33\x34\x35\x36\x37\x38\x39\x3A\x3B\x3C\x3D\x3E\x3F\x40\x41\x42\x43\x44\x45\x46\x47\x48\x49\x4A\x4B\x4C\x4D\x4E\x4F\x50\x51\x52\x53\x54\x55\x56\x57\x58\x59\x5A\x5B\x5C\x5D\x5E\x5F\x60\x61\x62\x63\x64\x65\x66\x67\x68\x69\x6A\x6B\x6C\x6D\x6E\x6F\x70\x71\x72\x73\x74\x75\x76\x77\x78\x79\x7A\x7B\x7C\x7D\x7E\x7F\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8A\x8B\x8C\x8D\x8E\x8F\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9A\x9B\x9C\x9D\x9E\x9F\xA0\xA1\xA2\xA3\xA4\xA5\xA6\xA7\xA8\xA9\xAA\xAB\xAC\xAD\xAE\xAF\xB0\xB1\xB2\xB3\xB4\xB5\xB6\xB7\xB8\xB9\xBA\xBB\xBC\xBD\xBE\xBF\xC0\xC1\xC2\xC3\xC4\xC5\xC6\xC7\xC8\xC9\xCA\xCB\xCC\xCD\xCE\xCF\xD0\xD1\xD2\xD3\xD4\xD5\xD6\xD7\xD8\xD9\xDA\xDB\xDC\xDD\xDE\xDF\xE0\xE1\xE2\xE3\xE4\xE5\xE6\xE7\xE8\xE9\xEA\xEB\xEC\xED\xEE\xEF\xF0\xF1\xF2\xF3\xF4\xF5\xF6\xF7\xF8\xF9\xFA\xFB\xFC\xFD\xFE\xFF", want: "$\"foo\\01\\02\\03\\04\\05\\06\\07\\08\\09\\0A\\0B\\0C\\0D\\0E\\0F\\10\\11\\12\\13\\14\\15\\16\\17\\18\\19\\1A\\1B\\1C\\1D\\1E\\1F !\\22#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\5C]^_`abcdefghijklmnopqrstuvwxyz{|}~\\7F\\80\\81\\82\\83\\84\\85\\86\\87\\88\\89\\8A\\8B\\8C\\8D\\8E\\8F\\90\\91\\92\\93\\94\\95\\96\\97\\98\\99\\9A\\9B\\9C\\9D\\9E\\9F\\A0\\A1\\A2\\A3\\A4\\A5\\A6\\A7\\A8\\A9\\AA\\AB\\AC\\AD\\AE\\AF\\B0\\B1\\B2\\B3\\B4\\B5\\B6\\B7\\B8\\B9\\BA\\BB\\BC\\BD\\BE\\BF\\C0\\C1\\C2\\C3\\C4\\C5\\C6\\C7\\C8\\C9\\CA\\CB\\CC\\CD\\CE\\CF\\D0\\D1\\D2\\D3\\D4\\D5\\D6\\D7\\D8\\D9\\DA\\DB\\DC\\DD\\DE\\DF\\E0\\E1\\E2\\E3\\E4\\E5\\E6\\E7\\E8\\E9\\EA\\EB\\EC\\ED\\EE\\EF\\F0\\F1\\F2\\F3\\F4\\F5\\F6\\F7\\F8\\F9\\FA\\FB\\FC\\FD\\FE\\FF\""},
	}
	for i, g := range golden {
		got := ComdatName(g.s)
		if g.want != got {
			t.Errorf("i=%d: name mismatch; expected `%s`, got `%s`", i, g.want, got)
		}
	}
}

func TestMetadataName(t *testing.T) {
	golden := []struct {
		s    string
		want string
	}{
		// i=0
		{s: "foo", want: "!foo"},
		// i=1
		{s: "a b", want: `!a\20b`},
		// i=2
		{s: "$a", want: "!$a"},
		// i=3
		{s: "-a", want: "!-a"},
		// i=4
		{s: ".a", want: "!.a"},
		// i=5
		{s: "_a", want: "!_a"},
		// i=6
		{s: "#a", want: `!\23a`},
		// i=7
		{s: "a b#c", want: `!a\20b\23c`},
		// i=8
		{s: "foo世bar", want: `!foo\E4\B8\96bar`},
		// i=9
		{s: "qux\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29\x2A\x2B\x2C\x2D\x2E\x2F\x30\x31\x32\x33\x34\x35\x36\x37\x38\x39\x3A\x3B\x3C\x3D\x3E\x3F\x40\x41\x42\x43\x44\x45\x46\x47\x48\x49\x4A\x4B\x4C\x4D\x4E\x4F\x50\x51\x52\x53\x54\x55\x56\x57\x58\x59\x5A\x5B\x5C\x5D\x5E\x5F\x60\x61\x62\x63\x64\x65\x66\x67\x68\x69\x6A\x6B\x6C\x6D\x6E\x6F\x70\x71\x72\x73\x74\x75\x76\x77\x78\x79\x7A\x7B\x7C\x7D\x7E\x7F\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8A\x8B\x8C\x8D\x8E\x8F\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9A\x9B\x9C\x9D\x9E\x9F\xA0\xA1\xA2\xA3\xA4\xA5\xA6\xA7\xA8\xA9\xAA\xAB\xAC\xAD\xAE\xAF\xB0\xB1\xB2\xB3\xB4\xB5\xB6\xB7\xB8\xB9\xBA\xBB\xBC\xBD\xBE\xBF\xC0\xC1\xC2\xC3\xC4\xC5\xC6\xC7\xC8\xC9\xCA\xCB\xCC\xCD\xCE\xCF\xD0\xD1\xD2\xD3\xD4\xD5\xD6\xD7\xD8\xD9\xDA\xDB\xDC\xDD\xDE\xDF\xE0\xE1\xE2\xE3\xE4\xE5\xE6\xE7\xE8\xE9\xEA\xEB\xEC\xED\xEE\xEF\xF0\xF1\xF2\xF3\xF4\xF5\xF6\xF7\xF8\xF9\xFA\xFB\xFC\xFD\xFE\xFF", want: `!qux\01\02\03\04\05\06\07\08\09\0A\0B\0C\0D\0E\0F\10\11\12\13\14\15\16\17\18\19\1A\1B\1C\1D\1E\1F\20\21\22\23$\25\26\27\28\29\2A\2B\2C-.\2F0123456789\3A\3B\3C\3D\3E\3F\40ABCDEFGHIJKLMNOPQRSTUVWXYZ\5B\5C\5D\5E_\60abcdefghijklmnopqrstuvwxyz\7B\7C\7D\7E\7F\80\81\82\83\84\85\86\87\88\89\8A\8B\8C\8D\8E\8F\90\91\92\93\94\95\96\97\98\99\9A\9B\9C\9D\9E\9F\A0\A1\A2\A3\A4\A5\A6\A7\A8\A9\AA\AB\AC\AD\AE\AF\B0\B1\B2\B3\B4\B5\B6\B7\B8\B9\BA\BB\BC\BD\BE\BF\C0\C1\C2\C3\C4\C5\C6\C7\C8\C9\CA\CB\CC\CD\CE\CF\D0\D1\D2\D3\D4\D5\D6\D7\D8\D9\DA\DB\DC\DD\DE\DF\E0\E1\E2\E3\E4\E5\E6\E7\E8\E9\EA\EB\EC\ED\EE\EF\F0\F1\F2\F3\F4\F5\F6\F7\F8\F9\FA\FB\FC\FD\FE\FF`},
		// i=10
		{s: "123", want: `!\3123`},
	}
	for i, g := range golden {
		got := MetadataName(g.s)
		if g.want != got {
			t.Errorf("i=%d: name mismatch; expected `%s`, got `%s`", i, g.want, got)
		}
	}
}

func TestMetadataID(t *testing.T) {
	golden := []struct {
		s    int64
		want string
	}{
		// i=0
		{s: 2, want: "!2"},
		// i=1
		{s: 42, want: "!42"},
	}
	for i, g := range golden {
		got := MetadataID(g.s)
		if g.want != got {
			t.Errorf("i=%d: name mismatch; expected `%s`, got `%s`", i, g.want, got)
		}
	}
}

func TestEscapeString(t *testing.T) {
	golden := []struct {
		s    []byte
		want string
	}{
		// i=0
		{s: []byte("foo"), want: "foo"},
		// i=1
		{s: []byte("a b"), want: `a b`},
		// i=2
		{s: []byte("$a"), want: "$a"},
		// i=3
		{s: []byte("-a"), want: "-a"},
		// i=4
		{s: []byte(".a"), want: ".a"},
		// i=5
		{s: []byte("_a"), want: "_a"},
		// i=6
		{s: []byte("#a"), want: `#a`},
		// i=7
		{s: []byte("a b#c"), want: `a b#c`},
		// i=8
		{s: []byte("2"), want: "2"},
		// i=9
		{s: []byte("foo世bar"), want: `foo\E4\B8\96bar`},
		// i=10
		{s: []byte(`foo \ bar`), want: `foo \5C bar`},
		// i=11 (arbitrary data, invalid UTF-8)
		{s: []byte{'f', 'o', 'o', 0x81, 0x82, 'b', 'a', 'r'}, want: `foo\81\82bar`},
	}
	for i, g := range golden {
		got := EscapeString(g.s)
		if g.want != got {
			t.Errorf("i=%d: string mismatch; expected `%s`, got `%s`", i, g.want, got)
		}
	}
}

func TestEscape(t *testing.T) {
	golden := []struct {
		s    []byte
		want string
	}{
		// i=0
		{s: []byte("foo"), want: "foo"},
		// i=1
		{s: []byte("a b"), want: `a b`},
		// i=2
		{s: []byte("$a"), want: "$a"},
		// i=3
		{s: []byte("-a"), want: "-a"},
		// i=4
		{s: []byte(".a"), want: ".a"},
		// i=5
		{s: []byte("_a"), want: "_a"},
		// i=6
		{s: []byte("#a"), want: `#a`},
		// i=7
		{s: []byte("a b#c"), want: `a b#c`},
		// i=8
		{s: []byte("2"), want: "2"},
		// i=9
		{s: []byte("foo世bar"), want: `foo\E4\B8\96bar`},
		// i=10
		{s: []byte(`foo \ bar`), want: `foo \5C bar`},
		// i=11 (arbitrary data, invalid UTF-8)
		{s: []byte{'f', 'o', 'o', 0x81, 0x82, 'b', 'a', 'r'}, want: `foo\81\82bar`},
	}
	// isPrint reports whether the given byte is printable in ASCII.
	isPrint := func(b byte) bool {
		return ' ' <= b && b <= '~' && b != '"' && b != '\\'
	}
	for i, g := range golden {
		got := Escape(g.s, isPrint)
		if g.want != got {
			t.Errorf("i=%d: string mismatch; expected `%s`, got `%s`", i, g.want, got)
		}
	}
}

func TestUnescape(t *testing.T) {
	golden := []struct {
		s    string
		want []byte
	}{
		// i=0
		{s: "foo", want: []byte("foo")},
		// i=1
		{s: `a\20b`, want: []byte("a b")},
		// i=2
		{s: "$a", want: []byte("$a")},
		// i=3
		{s: "-a", want: []byte("-a")},
		// i=4
		{s: ".a", want: []byte(".a")},
		// i=5
		{s: "_a", want: []byte("_a")},
		// i=6
		{s: `\23a`, want: []byte("#a")},
		// i=7
		{s: `a\20b\23c`, want: []byte("a b#c")},
		// i=8
		{s: "2", want: []byte("2")},
		// i=9
		{s: `foo\E4\B8\96bar`, want: []byte("foo世bar")},
		// i=10
		{s: `foo \5C bar`, want: []byte(`foo \ bar`)},
		// i=11
		{s: `foo \5c bar`, want: []byte(`foo \ bar`)},
		// i=12
		{s: `foo \\ bar`, want: []byte(`foo \ bar`)},
		// i=13 (arbitrary data, invalid UTF-8)
		{s: `foo\81\82bar`, want: []byte{'f', 'o', 'o', 0x81, 0x82, 'b', 'a', 'r'}},
	}
	for i, g := range golden {
		got := Unescape(g.s)
		if !reflect.DeepEqual(g.want, got) {
			t.Errorf("i=%d: string mismatch; expected `%s`, got `%s`", i, g.want, got)
		}
	}
}

func TestQuote(t *testing.T) {
	golden := []struct {
		s    []byte
		want string
	}{
		// i=0
		{s: []byte("foo"), want: `"foo"`},
		// i=1
		{s: []byte("a b"), want: `"a b"`},
		// i=2
		{s: []byte("$a"), want: `"$a"`},
		// i=3
		{s: []byte("-a"), want: `"-a"`},
		// i=4
		{s: []byte(".a"), want: `".a"`},
		// i=5
		{s: []byte("_a"), want: `"_a"`},
		// i=6
		{s: []byte("#a"), want: `"#a"`},
		// i=7
		{s: []byte("a b#c"), want: `"a b#c"`},
		// i=8
		{s: []byte("2"), want: `"2"`},
		// i=9
		{s: []byte("foo世bar"), want: `"foo\E4\B8\96bar"`},
		// i=10
		{s: []byte(`foo \ bar`), want: `"foo \5C bar"`},
		// i=11 (arbitrary data, invalid UTF-8)
		{s: []byte{'f', 'o', 'o', 0x81, 0x82, 'b', 'a', 'r'}, want: `"foo\81\82bar"`},
	}
	for i, g := range golden {
		got := Quote(g.s)
		if g.want != got {
			t.Errorf("i=%d: string mismatch; expected `%s`, got `%s`", i, g.want, got)
		}
	}
}

func TestUnquote(t *testing.T) {
	golden := []struct {
		s    string
		want []byte
	}{
		// i=0
		{s: `"foo"`, want: []byte("foo")},
		// i=1
		{s: `"a\20b"`, want: []byte("a b")},
		// i=2
		{s: `"$a"`, want: []byte("$a")},
		// i=3
		{s: `"-a"`, want: []byte("-a")},
		// i=4
		{s: `".a"`, want: []byte(".a")},
		// i=5
		{s: `"_a"`, want: []byte("_a")},
		// i=6
		{s: `"\23a"`, want: []byte("#a")},
		// i=7
		{s: `"a\20b\23c"`, want: []byte("a b#c")},
		// i=8
		{s: `"2"`, want: []byte("2")},
		// i=9
		{s: `"foo\E4\B8\96bar"`, want: []byte("foo世bar")},
		// i=10
		{s: `"foo \5C bar"`, want: []byte(`foo \ bar`)},
		// i=11
		{s: `"foo \\ bar"`, want: []byte(`foo \ bar`)},
		// i=12 (arbitrary data, invalid UTF-8)
		{s: `"foo\81\82bar"`, want: []byte{'f', 'o', 'o', 0x81, 0x82, 'b', 'a', 'r'}},
		// i=13
		{s: `"\8E\\"`, want: []byte{0x8E, '\\'}}, // see issue llir/llvm#184
		// i=14
		{s: `"\\"`, want: []byte{'\\'}},
	}
	for i, g := range golden {
		got := Unquote(g.s)
		if !reflect.DeepEqual(g.want, got) {
			t.Errorf("i=%d: string mismatch; expected `%s`, got `%s`", i, g.want, got)
		}
	}
}

func BenchmarkGlobalNoReplace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GlobalName("$foo_bar_baz")
	}
}

func BenchmarkGlobalReplace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GlobalName("$foo bar#baz")
	}
}

func BenchmarkLocalNoReplace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LocalName("$foo_bar_baz")
	}
}

func BenchmarkLocalReplace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LocalName("$foo bar#baz")
	}
}
