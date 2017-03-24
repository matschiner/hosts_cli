package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAddition(t *testing.T) {
	ip := "1.1.1.1"
	hn := "example.com"

	Convey("given a hosts file", t, func() {
		f := &hostlist{}

		Convey("when a host is added", func() {
			So(f.Add(ip, hn), ShouldBeNil)

			Convey("it should appear in the list", func() {
				So(f.lines, ShouldContain, ip+"\t"+hn)
			})

			Convey("it should be changed", func() {
				So(f.changed, ShouldBeTrue)
			})

			Convey("should be in the bytes", func() {
				So(f.Bytes(), ShouldResemble, []byte(ip+"\t"+hn))
			})

		})
	})
}

func TestRemoval(t *testing.T) {
	ip := "1.1.1.1"
	hn := "example.com"

	Convey("given a hosts file", t, func() {
		f := &hostlist{}

		Convey("with an entry", func() {
			f.lines = []string{ip + "\t" + hn}

			Convey("when it is removed", func() {
				Convey("via the ip", func() {
					So(f.Remove(ip), ShouldBeNil)

					Convey("it should not be in the lines", func() {
						So(f.lines, ShouldNotContain, ip+"\t"+hn)
					})
				})

				Convey("via the hostname", func() {
					So(f.Remove(hn), ShouldBeNil)

					Convey("it should not be in the lines", func() {
						So(f.lines, ShouldNotContain, ip+"\t"+hn)
					})
				})

				Convey("via part of the IP", func() {
					So(f.Remove("1.1."), ShouldBeNil)
					Convey("it should NOT be removed", func() {
						So(f.lines, ShouldContain, ip+"\t"+hn)
					})
				})

				Convey("via part of the hostname", func() {
					So(f.Remove("mple.co"), ShouldBeNil)
					Convey("it should NOT be removed", func() {
						So(f.lines, ShouldContain, ip+"\t"+hn)
					})
				})
			})

		})
	})
}

func TestContains(t *testing.T) {
	ip := "1.1.1.1"
	hn := "example.com"

	Convey("given a hosts file", t, func() {
		f := &hostlist{}

		Convey("with an entry", func() {
			f.lines = []string{ip + "\t" + hn}

			Convey("it should know that entry is contained", func() {
				ok, err := f.Contains(ip, hn)
				So(err, ShouldBeNil)
				So(ok, ShouldBeTrue)

				ok, err = f.Contains(hn, ip)
				So(err, ShouldBeNil)
				So(ok, ShouldBeTrue)
			})
		})
	})
}

func TestParsing(t *testing.T) {
	ip := "1.1.1.1"
	ip2 := "1.1.1.2"
	hn := "example.com"
	hn2 := "www.example.com"

	Convey("given a hosts file", t, func() {
		f := &hostlist{}

		Convey("when parsing incoming bytes", func() {
			Convey("it should add the line to the entries", func() {
				f.Parse([]byte(ip + "\t" + hn))
				So(f.lines, ShouldContain, ip+"\t"+hn)
			})

			Convey("it should condense multiple new lines", func() {
				f.Parse([]byte(ip + "\t" + hn + "\n\n\n\n\n\n" + ip2 + "\t" + hn2))
				So(f.Bytes(), ShouldResemble, []byte(ip+"\t"+hn+"\n\n"+ip2+"\t"+hn2))
			})
		})
	})
}

func TestCommenting(t *testing.T) {
	ip := "1.1.1.1"
	hn := "example.com"

	Convey("given a hosts file", t, func() {
		f := &hostlist{}

		Convey("with an entry", func() {
			f.lines = []string{ip + "\t" + hn}

			Convey("should be able to be commented", func() {
				Convey("via the hostname", func() {
					So(f.Comment(hn), ShouldBeNil)
					So(f.lines, ShouldNotContain, ip+"\t"+hn)
					So(f.lines, ShouldContain, "#"+ip+"\t"+hn)

					Convey("and uncommented", func() {
						Convey("via the hostname", func() {
							So(f.Uncomment(hn), ShouldBeNil)
							So(f.lines, ShouldContain, ip+"\t"+hn)
							So(f.lines, ShouldNotContain, "#"+ip+"\t"+hn)
						})

						Convey("via the ip", func() {
							So(f.Uncomment(ip), ShouldBeNil)
							So(f.lines, ShouldNotContain, "#"+ip+"\t"+hn)
							So(f.lines, ShouldContain, ip+"\t"+hn)
						})
					})
				})

				Convey("via the ip", func() {
					So(f.Comment(ip), ShouldBeNil)
					So(f.lines, ShouldNotContain, ip+"\t"+hn)
					So(f.lines, ShouldContain, "#"+ip+"\t"+hn)

					Convey("and uncommented", func() {
						Convey("via the hostname", func() {
							So(f.Uncomment(hn), ShouldBeNil)
							So(f.lines, ShouldContain, ip+"\t"+hn)
							So(f.lines, ShouldNotContain, "#"+ip+"\t"+hn)
						})

						Convey("via the ip", func() {
							So(f.Uncomment(ip), ShouldBeNil)
							So(f.lines, ShouldNotContain, "#"+ip+"\t"+hn)
							So(f.lines, ShouldContain, ip+"\t"+hn)
						})
					})
				})
			})

		})
	})
}
