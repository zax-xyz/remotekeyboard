C := go build
FLAGS := -ldflags "-s -w"

TARGET := remotekeyboard

all: ${TARGET}

${TARGET}: ${TARGET}.go
	${C} ${FLAGS} $@.go
	upx $@

clean:
	rm ${TARGET}
