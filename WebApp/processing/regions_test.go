package processing

import (
	"compress/gzip"
	"encoding/binary"
	"errors"
	"flag"
	"image"
	"image/png"
	"io"
	"math/bits"
	"os"
	"slices"
	"testing"
)

var (
	updateSnapshots = flag.Bool("update-snapshots", false, "Update snapshots to reflect the current behavior of the algorithm.")
)

func loadImage(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	i, err := png.Decode(f)
	if err != nil {
		panic(err)
	}

	return i
}

type regionTestArgs struct {
	img          image.Image
	options      RegionMapOptions
	regionFilter func(Region) bool
}

type regionTest struct {
	name string
	args regionTestArgs
}

var regionTests = []regionTest{
	{
		name: "lub",
		args: regionTestArgs{
			img: loadImage("./build_region_map_test_images/lub.png"),
		},
	},
	{
		name: "whiteboardshapes",
		args: regionTestArgs{
			img: loadImage("./build_region_map_test_images/whiteboardshapes.png"),
		},
	},
	{
		name: "allwhite",
		args: regionTestArgs{
			img: loadImage("./build_region_map_test_images/allwhite.png"),
		},
	},
	{
		name: "allblack",
		args: regionTestArgs{
			img: loadImage("./build_region_map_test_images/allblack.png"),
		},
	},
	{
		name: "allcolors",
		args: regionTestArgs{
			img: loadImage("./build_region_map_test_images/allcolors.png"),
		},
	},
}

func nearestIntBitLength(len int) int {
	switch {
	case len <= 8:
		return 8
	case len <= 16:
		return 16
	case len <= 32:
		return 32
	case len <= 64:
		return 64
	default:
		panic("TOO BIG!!!")
	}
}

func TestBuildRegionMap(t *testing.T) {
	flag.Parse()
	for _, tt := range regionTests {
		t.Run(tt.name, func(t *testing.T) {
			img := SimplifyImage(tt.args.img, tt.args.options)
			rm := BuildRegionMap(img, tt.args.options, tt.args.regionFilter)

			pixelCount := img.Bounds().Dx() * img.Bounds().Dy()

			currentRegionIds := make([]uint64, pixelCount)

			stride := img.Bounds().Dx()

			sortedRegions := slices.SortedFunc(slices.Values(rm.GetRegions()), func(a, b Region) int {
				return (a.GetBounds().Min.X + a.GetBounds().Min.Y*stride) - (b.GetBounds().Min.X + b.GetBounds().Min.Y*stride)
			})

			for i, v := range sortedRegions {
				for _, p := range v {
					currentRegionIds[int(p.X)+int(p.Y)*stride] = uint64(i + 1)
				}
			}

			if !*updateSnapshots {
				snapFile, err := os.Open("./snapshots/build_region_map_" + tt.name + ".snap")
				if errors.Is(err, os.ErrNotExist) {
					t.Fatalf("Snapshot not created for test %s; Run the test with the -update-snapshots flag.", tt.name)
				}
				if err != nil {
					panic(err)
				}
				defer snapFile.Close()
				gzipReader, err := gzip.NewReader(snapFile)
				if err != nil {
					panic(err)
				}
				defer gzipReader.Close()

				var regionIdBitLength int
				{
					b := make([]byte, 1)
					_, err = gzipReader.Read(b)
					if err != nil {
						panic(err)
					}
					regionIdBitLength = int(b[0])
				}

				if bl := nearestIntBitLength(bits.Len(uint(len(sortedRegions)))); bl != regionIdBitLength {
					t.Fatalf("Snapshot %s failed: region id bit length mismatch (%d != %d)", tt.name, bl, regionIdBitLength)
				}

				targetRegionIdsBytes, err := io.ReadAll(gzipReader)
				if err != nil {
					panic(err)
				}

				switch regionIdBitLength {
				case 8:
					if len(targetRegionIdsBytes) != img.Bounds().Dx()*img.Bounds().Dy() {
						t.Fatalf("Snapshot %s failed: snapshot regionmap size inconsistent with current image", tt.name)
					}

					if len(currentRegionIds) != len(targetRegionIdsBytes) {
						t.Fatalf("Snapshot %s failed: current regionmap SIZE does not match snapshot", tt.name)
					}

					for i, v := range currentRegionIds {
						if byte(v) != targetRegionIdsBytes[i] {
							t.Fatalf("Snapshot %s failed: current regionmap does not match snapshot", tt.name)
						}
					}

				case 16:
					if len(targetRegionIdsBytes) != img.Bounds().Dx()*img.Bounds().Dy()*2 {
						t.Fatalf("Snapshot %s failed: inconsistent length", tt.name)
					}
					targetRegionIds := make([]uint16, len(targetRegionIdsBytes)/2)

					binary.Decode(targetRegionIdsBytes, binary.LittleEndian, targetRegionIds)

					if len(currentRegionIds) != len(targetRegionIds) {
						t.Fatalf("Snapshot %s failed: current regionmap SIZE does not match snapshot", tt.name)
					}

					for i, v := range currentRegionIds {
						if uint16(v) != targetRegionIds[i] {
							t.Fatalf("Snapshot %s failed: current regionmap does not match snapshot", tt.name)
						}
					}
				case 32:
					if len(targetRegionIdsBytes) != img.Bounds().Dx()*img.Bounds().Dy()*4 {
						t.Fatalf("Snapshot %s failed: inconsistent length", tt.name)
					}
					targetRegionIds := make([]uint32, len(targetRegionIdsBytes)/2)

					binary.Decode(targetRegionIdsBytes, binary.LittleEndian, targetRegionIds)

					if len(currentRegionIds) != len(targetRegionIds) {
						t.Fatalf("Snapshot %s failed: current regionmap SIZE does not match snapshot", tt.name)
					}

					for i, v := range currentRegionIds {
						if uint32(v) != targetRegionIds[i] {
							t.Fatalf("Snapshot %s failed: current regionmap does not match snapshot", tt.name)
						}
					}
				case 64:
					if len(targetRegionIdsBytes) != img.Bounds().Dx()*img.Bounds().Dy()*8 {
						t.Fatalf("Snapshot %s failed: inconsistent length", tt.name)
					}
					targetRegionIds := make([]uint64, len(targetRegionIdsBytes)/2)

					binary.Decode(targetRegionIdsBytes, binary.LittleEndian, targetRegionIds)

					if len(currentRegionIds) != len(targetRegionIds) {
						t.Fatalf("Snapshot %s failed: current regionmap SIZE does not match snapshot", tt.name)
					}

					for i, v := range currentRegionIds {
						if v != targetRegionIds[i] {
							t.Fatalf("Snapshot %s failed: current regionmap does not match snapshot", tt.name)
						}
					}
				default:
					t.Fatalf("Snapshot %s failed: snapshot file is weird/corrupt, update the snapshot.", tt.name)
				}
			} else {
				err := os.MkdirAll("snapshots", os.ModeAppend)
				if err != nil {
					panic(err)
				}
				snapFile, err := os.Create("./snapshots/build_region_map_" + tt.name + ".snap")
				if err != nil {
					panic(err)
				}
				defer snapFile.Close()

				gzipWriter := gzip.NewWriter(snapFile)
				defer gzipWriter.Close()

				bitsRequired := bits.Len(uint(len(sortedRegions)))
				switch {
				case bitsRequired <= 8:
					gzipWriter.Write([]byte{8})
					dat := make([]byte, len(currentRegionIds))
					for i, v := range currentRegionIds {
						dat[i] = byte(v)
					}
					gzipWriter.Write(dat)
				case bitsRequired <= 16:
					gzipWriter.Write([]byte{16})
					dat := make([]uint16, len(currentRegionIds))
					for i, v := range currentRegionIds {
						dat[i] = uint16(v)
					}
					binary.Write(gzipWriter, binary.LittleEndian, dat)
				case bitsRequired <= 32:
					gzipWriter.Write([]byte{32})
					dat := make([]uint32, len(currentRegionIds))
					for i, v := range currentRegionIds {
						dat[i] = uint32(v)
					}
					binary.Write(gzipWriter, binary.LittleEndian, dat)
				case bitsRequired <= 64:
					gzipWriter.Write([]byte{64})
					binary.Write(gzipWriter, binary.LittleEndian, currentRegionIds)
				default:
					panic("What just happened, what are you doing dude")
				}
			}
		})
	}
}

func BenchmarkBuildRegionMap(b *testing.B) {
	for _, bm := range regionTests {
		b.Run(bm.name, func(b *testing.B) {
			img := SimplifyImage(bm.args.img, bm.args.options)
			for b.Loop() {
				BuildRegionMap(img, bm.args.options, bm.args.regionFilter)
			}
		})
	}
}
