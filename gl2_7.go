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

	rendering_program := compileProgram()

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var currentTime float64 = 0
	var currTime time.Time = time.Now()
	for !window.ShouldClose() {
		render(currentTime, rendering_program)
		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()

		nextTime := time.Now()
		currentTime += nextTime.Sub(currTime).Seconds()
		currTime = nextTime
	}

	gl.DeleteProgram(rendering_program)
	gl.DeleteVertexArrays(1, &vao)
}

func render(currentTime float64, rendering_program uint32) {
	red := []float32{
		float32(math.Sin(currentTime))*0.5 + 0.5,
		float32(math.Cos(currentTime))*0.5 + 0.5,
		0.0, 1.0}
	gl.ClearBufferfv(gl.COLOR, 0, &red[0])

	gl.UseProgram(rendering_program)

	gl.DrawArrays(gl.POINTS, 0, 1)
	gl.PointSize(40)
}

// 간단한 쉐이더 컴파일하기
func compileProgram() uint32 {
	var vertex_shader_source = `
		#version 430 core

		void main(void) {
			gl_Position = vec4(0.0, 0.0, 0.5, 1.0);
		}
	` + "\x00"

	var fragment_shader_source = `
		#version 430 core

		out vec4 color;

		void main(void) {
			color = vec4(0.0, 0.8, 1.0, 1.0);
		}
	` + "\x00"

	vertex_shader := gl.CreateShader(gl.VERTEX_SHADER)
	gl_vertex_shader_source := gl.Str(vertex_shader_source)
	gl.ShaderSource(vertex_shader, 1, &gl_vertex_shader_source, nil)
	gl.CompileShader(vertex_shader)

	fragment_shader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl_fragment_shader_source := gl.Str(fragment_shader_source)
	gl.ShaderSource(fragment_shader, 1, &gl_fragment_shader_source, nil)
	gl.CompileShader(fragment_shader)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertex_shader)
	gl.AttachShader(program, fragment_shader)
	gl.LinkProgram(program)

	gl.DeleteShader(vertex_shader)
	gl.DeleteShader(fragment_shader)

	return program
}
