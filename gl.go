package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const WindowWidth = 800
const WindowHeight = 600

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	// glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	// glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Test", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	vendor := gl.GoStr(gl.GetString(gl.VENDOR))
	renderer := gl.GoStr(gl.GetString(gl.RENDERER))
	fmt.Println("OpenGL version", version)
	fmt.Println("OpenGL vendor", vendor)
	fmt.Println("OpenGL renderer", renderer)

	var currentTime float64 = 0
	var currTime time.Time = time.Now()
	for !window.ShouldClose() {
		render(currentTime)
		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()

		nextTime := time.Now()
		currentTime += nextTime.Sub(currTime).Seconds()
		currTime = nextTime
	}
}

func render(currentTime float64) {
	// example2_1(currentTime)
	// example2_2(currentTime)
	example2_3(currentTime)
}

// 간단한 쉐이더 컴파일하기
func example2_3(currentTime float64) {

	var vertex_shader_source = `
		#version 430 core

		void main(void) {
			gl_Position = vec4(0.0, 0.0, 0.5, 1.0);
		}
	`

	var fragment_shader_source = `
		#version 430 core

		out vec4 color;

		void main(void) {
			color = vec4(0.0, 0.8, 1.0, 1.0)
		}
	`
	vertex_shader := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertexShader, 1, vertex_shader_source, nil)
	gl.CompileShader(vertexShader)

	fragment_shader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragment_shader, 1, fragment_shader_source, nil)
	gl.CompileShader(fragment_shader)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertex_shader)
}

// 시간에 따라 color 변화시키기
func example2_2(currentTime float64) {
	red := []float32{
		float32(math.Sin(currentTime))*0.5 + 0.5,
		float32(math.Cos(currentTime))*0.5 + 0.5,
		0.0, 1.0}
	gl.ClearBufferfv(gl.COLOR, 0, &red[0])
}

// 첫 번째 OpenGL 애플리케이션
func example2_1(currentTime float64) {
	red := []float32{0.0, 0.0, 1.0, 1.0}
	gl.ClearBufferfv(gl.COLOR, 0, &red[0])
}
