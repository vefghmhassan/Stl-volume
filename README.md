
# STLCalc: STL Volume and Weight Calculator
STLCalc is a Go library designed to calculate the volume and weight of 3D objects described by STL (Stereolithography) files. It supports both binary and ASCII STL file formats, allowing users to easily analyze and manipulate 3D object data, particularly useful in 3D printing, CAD software, or similar applications.

Features
File Format Support: Handles both ASCII and binary STL files.
Volume Calculation: Computes the volume of STL objects.
Weight Calculation: Calculates the object's weight based on volume and material density.
Material Density: Allows setting custom material densities.
Triangle Count: Retrieves the number of triangles in the STL file.
Installation
To use STLCalc, you'll need to have Go installed on your system. You can then include STLCalc in your Go project by adding the provided stlcalc.go file to your project directory.

Usage
Initializing STLCalc
First, initialize a new instance of STLCalc by providing the path to your STL file:
```go
calc, err := NewSTLCalc("path/to/your/object.stl")
if err != nil {
    log.Fatalf("Failed to initialize STLCalc: %v", err)
}
defer calc.Close()
```
# Setting Material Density
Set the density of the material (in g/cm³) for accurate weight calculations:
```go
calc.SetDensity(1.04) // Example for ABS plastic
```

# Calculating Volume and Weight
Calculate the volume (in cm³ or in³) and weight (in grams) of the object:
```go
volumeCm3, err := calc.GetVolume("cm") // or "in" for cubic inches
if err != nil {
    log.Fatalf("Failed to calculate volume: %v", err)
}

weight, err := calc.GetWeight()
if err != nil {
    log.Fatalf("Failed to calculate weight: %v", err)
}
```
# Retrieving Triangle Count
Optionally, get the number of triangles in the STL file:

```go
trianglesCount := calc.GetTrianglesCount()
```


