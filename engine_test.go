package lua_test

import (
	. "github.com/seer-server/script-engine"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Engine", func() {
	var (
		err          error
		engine       *Engine
		fileName     = "lua_test.lua"
		stringScript = `
			function hello(name)
				return "Hello, " .. name .. "!"
			end
		`
	)

	BeforeEach(func() {
		engine = NewEngine()
	})

	Context("when loading from a string", func() {
		BeforeEach(func() {
			err = engine.LoadString(stringScript)
		})

		It("should not fail", func() {
			Expect(err).To(BeNil())
		})

		It("should be able to call a method", func() {
			results, err := engine.Call("hello", 1, "World")
			Expect(err).To(BeNil())
			Expect(len(results)).To(Equal(1))
			Expect(results[0]).ToNot(Equal(Nil))
			Expect(results[0].AsString()).To(Equal("Hello, World!"))
		})
	})

	Context("when loading from a file", func() {
		BeforeEach(func() {
			err = engine.LoadFile(fileName)
		})

		It("shoult not fail", func() {
			Expect(err).To(BeNil())
		})

		It("should be able to call a method", func() {
			results, err := engine.Call("give_me_one", 1)
			Expect(err).To(BeNil())
			Expect(len(results)).To(Equal(1))
			Expect(results[0]).NotTo(Equal(Nil))
			Expect(results[0].AsNumber()).To(Equal(float64(1)))
		})
	})
})
