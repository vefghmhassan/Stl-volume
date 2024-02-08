package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"regexp"
)

type STLCalc struct {
	volume         float64
	weight         float64
	density        float64
	trianglesCount uint32
	bBinary        bool
	fstlHandle     *os.File
	fstlPath       string
	flag           bool
}

// NewSTLCalc initializes the STLCalc struct
func NewSTLCalc(filepath string) (*STLCalc, error) {
	calc := &STLCalc{
		density: 1.04,
	}

	b, err := isAscii(filepath)
	if err != nil {
		return nil, err
	}

	if !b {
		fmt.Println("BINARY STL Suspected.")
		calc.bBinary = true
		file, err := os.Open(filepath)
		if err != nil {
			return nil, err
		}
		calc.fstlHandle = file
		calc.fstlPath = filepath
	} else {
		fmt.Println("ASCII STL Suspected.")
	}

	return calc, nil
}

// Close cleans up any open resources
func (c *STLCalc) Close() {
	if c.fstlHandle != nil {
		c.fstlHandle.Close()
	}
}

// GetVolume calculates and returns the volume of the STL object
func (c *STLCalc) GetVolume(unit string) (float64, error) {
	if !c.flag {
		volume, err := c.calculateVolume()
		if err != nil {
			return 0, err
		}
		c.volume = volume
		c.flag = true
	}
	if unit == "cm" {
		return c.volume / 1000, nil
	}
	return c.inch3(c.volume / 1000), nil
}

// GetWeight calculates and returns the weight of the STL object
func (c *STLCalc) GetWeight() (float64, error) {
	volume, err := c.GetVolume("cm")
	if err != nil {
		return 0, err
	}
	return c.calculateWeight(volume), nil
}

// SetDensity sets the density of the material
func (c *STLCalc) SetDensity(density float64) {
	c.density = density
}

// GetDensity returns the density of the material
func (c *STLCalc) GetDensity() float64 {
	return c.density
}

// GetTrianglesCount returns the number of triangles in the STL file
func (c *STLCalc) GetTrianglesCount() uint32 {
	return c.trianglesCount
}

// calculateVolume calculates the volume of the STL object
func (c *STLCalc) calculateVolume() (float64, error) {
	if c.bBinary {
		return c.readBinaryVolume()
	}
	// Implement ASCII STL volume calculation if needed
	return 0, nil
}

// readBinaryVolume reads and calculates the volume from a binary STL file
func (c *STLCalc) readBinaryVolume() (float64, error) {
	// Skip the header
	_, err := c.fstlHandle.Seek(80, 0)
	if err != nil {
		return 0, err
	}

	// Read the number of triangles
	var count uint32
	err = binary.Read(c.fstlHandle, binary.LittleEndian, &count)
	if err != nil {
		return 0, err
	}
	c.trianglesCount = count

	totalVolume := 0.0
	for i := uint32(0); i < count; i++ {
		volume, err := c.readTriangle()
		if err != nil {
			return 0, err
		}
		totalVolume += volume
	}

	return math.Abs(totalVolume), nil
}

// readTriangle reads and calculates the volume of a single triangle from a binary STL file
func (c *STLCalc) readTriangle() (float64, error) {
	// Define a struct for reading the triangle data
	var normal [3]float32
	var vertex1 [3]float32
	var vertex2 [3]float32
	var vertex3 [3]float32
	var attrByteCount uint16

	// Read the data
	if err := binary.Read(c.fstlHandle, binary.LittleEndian, &normal); err != nil {
		return 0, err
	}
	if err := binary.Read(c.fstlHandle, binary.LittleEndian, &vertex1); err != nil {
		return 0, err
	}
	if err := binary.Read(c.fstlHandle, binary.LittleEndian, &vertex2); err != nil {
		return 0, err
	}
	if err := binary.Read(c.fstlHandle, binary.LittleEndian, &vertex3); err != nil {
		return 0, err
	}
	if err := binary.Read(c.fstlHandle, binary.LittleEndian, &attrByteCount); err != nil {
		return 0, err
	}

	volume := signedVolumeOfTriangle(vertex1, vertex2, vertex3)
	return volume, nil
}

// signedVolumeOfTriangle calculates the signed volume of a triangle
func signedVolumeOfTriangle(p1, p2, p3 [3]float32) float64 {
	v321 := float64(p3[0]) * float64(p2[1]) * float64(p1[2])
	v231 := float64(p2[0]) * float64(p3[1]) * float64(p1[2])
	v312 := float64(p3[0]) * float64(p1[1]) * float64(p2[2])
	v132 := float64(p1[0]) * float64(p3[1]) * float64(p2[2])
	v213 := float64(p2[0]) * float64(p1[1]) * float64(p3[2])
	v123 := float64(p1[0]) * float64(p2[1]) * float64(p3[2])
	return (1.0 / 6.0) * (-v321 + v231 + v312 - v132 - v213 + v123)
}

// inch3 converts cubic centimeters to cubic inches
func (c *STLCalc) inch3(volumeCm3 float64) float64 {
	return volumeCm3 * 0.0610237441
}

// calculateWeight calculates the weight of the object based on its volume and density
func (c *STLCalc) calculateWeight(volumeCm3 float64) float64 {
	return volumeCm3 * c.density
}

// isAscii checks if the given file is an ASCII STL file
func isAscii(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return false, err
	}

	asciiPattern := `facet normal\s+([-\d.eE]+)\s+([-\d.eE]+)\s+([-\d.eE]+)\s+outer loop\s+vertex\s+([-\d.eE]+)\s+([-\d.eE]+)\s+([-\d.eE]+)\s+vertex\s+([-\d.eE]+)\s+([-\d.eE]+)\s+([-\d.eE]+)\s+vertex\s+([-\d.eE]+)\s+([-\d.eE]+)\s+([-\d.eE]+)\s+endloop\s+endfacet`
	matched, err := regexp.Match(asciiPattern, buffer)
	if err != nil {
		return false, err
	}

	return matched, nil
}
