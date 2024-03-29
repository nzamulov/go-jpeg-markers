NAME 	 = gjm-cli
BIN_PATH = $(CURDIR)/$(NAME)

$(NAME):
	@go build -o $(BIN_PATH) $(CURDIR)/cli

clean:
	rm -f $(BIN_PATH)

.SILENT:
scan: clean $(NAME)
	$(BIN_PATH) -p $(lastword $(MAKECMDGOALS))

%.jpeg %.jpg:
	exit
