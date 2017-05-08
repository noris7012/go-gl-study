package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const windowWidth = 800
const windowHeight = 800

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Test", nil, nil)
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

	renderingProgram := compileProgram()

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var currentTime float64
	var currTime = time.Now()

	for !window.ShouldClose() {
		render(currentTime, renderingProgram)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()

		// <-time.After(10 * time.Millisecond)

		nextTime := time.Now()
		currentTime += nextTime.Sub(currTime).Seconds()
		currTime = nextTime
	}

	gl.DeleteProgram(renderingProgram)
	gl.DeleteVertexArrays(1, &vao)
}

func render(currentTime float64, renderingProgram uint32) {
	// 배경화면 색 조절
	color := []float32{
		// float32(math.Sin(currentTime))*0.5 + 0.5,
		// float32(math.Cos(currentTime))*0.5 + 0.5,
		0.0, 0.0, 0.0, 0.0}
	gl.ClearBufferfv(gl.COLOR, 0, &color[0])

	// shader 실행
	gl.UseProgram(renderingProgram)

	// locatin = 0, offset 과 대응
	attrib := []float32{
		float32(math.Sin(currentTime)) * 0.5,
		float32(math.Cos(currentTime)) * 0.05,
		0.0, 0.0}

	gl.VertexAttrib4fv(0, &attrib[0])

	// lcation = 1, color 와 대응
	color2 := []float32{
		// float32(math.Cos(currentTime))*0.5 + 0.5,
		float32(currentTime),
		float32(math.Sin(currentTime))*0.5 + 0.5,
		0.0, 1.0}

	gl.VertexAttrib4fv(1, &color2[0])

	gl.BeginTransformFeedback(gl.TRIANGLES)
	gl.EndTransformFeedback()
	// gl.DrawArrays(gl.TRIANGLES, 3, 3)
	// gl.DrawArrays(gl.TRIANGLES, 4, 3)
}

// 간단한 쉐이더 컴파일하기
func compileProgram() uint32 {
	var vertexShaderSource = `
		#version 430 core

		layout (location = 0) in vec4 offset;
		layout (location = 1) in vec4 color;

		out VS_OUT {
			vec4 color;
		} vs_color;

		void main(void) {
			const vec4 vertices[3] = vec4[3](
				vec4(0.25, -0.25, 0.5, 1.0),
				vec4(-0.25, -0.25, 0.5, 1.0),
				vec4(0, 0.25, 0.5, 1.0)
			);

			const vec4 vertices2[4] = vec4[4](
				vec4(-0.5, 0.5, 0.5, 1.0),
				vec4(0.5, 0.5, 0.5, 1.0),
				vec4(-0.5, -0.5, 0.5, 1.0),
				vec4(0.5, -0.5, 0.5, 1.0)
			);

			if (gl_VertexID < 3)
				gl_Position = vertices[gl_VertexID] + offset;
			else
				gl_Position = vertices2[gl_VertexID - 3] + offset;

      		vs_color.color = color;
		}
	` + "\x00"

	var fragmentShaderSource = `
		#version 430 core

		in VS_OUT {
			vec4 color;
		} fs_in;

		out vec4 color;

		void main(void) {
			color = fs_in.color;
		}
	` + "\x00"

	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	glVertexShaderSource := gl.Str(vertexShaderSource)
	gl.ShaderSource(vertexShader, 1, &glVertexShaderSource, nil)
	gl.CompileShader(vertexShader)

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	glFragmentShaderSource := gl.Str(fragmentShaderSource)
	gl.ShaderSource(fragmentShader, 1, &glFragmentShaderSource, nil)
	gl.CompileShader(fragmentShader)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program
}
