NAME 	 = gjm-cli
BIN_PATH = $(CURDIR)/$(NAME)

$(NAME):
	@go build -o $(BIN_PATH) $(CURDIR)

clean:
	rm -f $(BIN_PATH)