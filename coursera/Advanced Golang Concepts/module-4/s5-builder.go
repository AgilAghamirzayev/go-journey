package main

import "fmt"

type Computer struct {
	CPU     string
	GPU     string
	RAM     int
	Storage int
}

type Builder interface {
	SetCPU(cpu string) Builder
	SetGPU(gpu string) Builder
	SetRAM(ram int) Builder
	SetStorage(storage int) Builder
	Build() Computer
}

type ComputerBuilder struct {
	computer Computer
}

func NewComputerBuilder() *ComputerBuilder {
	return &ComputerBuilder{}
}

func (b *ComputerBuilder) SetCPU(cpu string) Builder {
	b.computer.CPU = cpu
	return b
}

func (b *ComputerBuilder) SetGPU(gpu string) Builder {
	b.computer.GPU = gpu
	return b
}

func (b *ComputerBuilder) SetRAM(ram int) Builder {
	b.computer.RAM = ram
	return b
}

func (b *ComputerBuilder) SetStorage(storage int) Builder {
	b.computer.Storage = storage
	return b
}

func (b *ComputerBuilder) Build() Computer {
	return b.computer
}

func main() {
	builder := NewComputerBuilder()

	gamingPC := builder.
		SetCPU("Intel i9").
		SetGPU("NVIDIA RTX 4090").
		SetRAM(32).
		SetStorage(2000).
		Build()

	officePC := builder.
		SetCPU("Intel i5").
		SetGPU("Integrated").
		SetRAM(16).
		SetStorage(512).
		Build()

	fmt.Printf("Gaming PC: %+v\n", gamingPC)
	fmt.Printf("Office PC: %+v\n", officePC)
}
