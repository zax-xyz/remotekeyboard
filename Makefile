C := go build
FLAGS := -ldflags "-s -w"

TARGET := remotekeyboard

all: ${TARGET}

${TARGET}:
	${C} ${FLAGS} $@.go
	upx $@
