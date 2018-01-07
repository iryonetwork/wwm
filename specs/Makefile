SOURCE_FILES = $(wildcard *.proto)
BUILT_FILES = $(subst .proto,.pb.go,$(SOURCE_FILES))
all: $(BUILT_FILES)

$(BUILT_FILES): $(SOURCE_FILES)
ifeq ("$(CI)", "")
	docker run -it --rm -v $$PWD:/tmp trifs/protoc:3.2.0 -I=/tmp --go_out=/tmp $(addprefix /tmp/,$(SOURCE_FILES))
else
	protoc -I=./ --go_out=. $(SOURCE_FILES)
endif

clear:
	rm -f *.pb.go
