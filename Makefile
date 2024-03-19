NAME 	 = gjm-cli
GO 		 = go
BIN_PATH = $(CURDIR)/$(NAME)

$(NAME):
	@$(GO) build -o $(BIN_PATH) $(CURDIR)

clean:
	rm -f $(BIN_PATH)