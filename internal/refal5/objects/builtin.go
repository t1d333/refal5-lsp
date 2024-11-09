package objects

var BuiltInFunctions = []RefalFunction{
	{
		Name:        "Card",
		Description: "Returns (is replaced by) the next line from the input. Normally it is from the terminal, but input can be redirected as allowed by MS-DOS. The returned expression is a sequence of character-symbols (possibly empty). The End-Of-Line byte is not included. If the input is read from a file, the macrodigit 0 is returned when the end of file is reached (no more lines). This is used in programs as the indicator of end, since a macrodigit cannot result from input otherwise. When reading from the terminal, enter Control-Z to produce the same effect",
		Category:    "I/O Functions",
		Signature:   "<Card>",
		ReturnValue: "Object",
	},
	{
		Name:        "Close ",
		Description: "s.ID is the identifying number of the file. The function closes the file. If such a file does not exist, no action is taken. The value of Close is always empty",
		Category:    "I/O Functions",
		Signature:   "<Close $0>",
	},
	{
		Name:        "ExistFile",
		Description: "e.name is a file name (a string of characters). The function checks whether the file exists and returnes the corresponding identifier \"True\" or \"False\"",
		Category:    "I/O Functions",
		Signature:   "",
		ReturnValue: "Boolean",
	},
	{
		Name:        "Print",
		Description: "prints the expression e.Expr on the current output and returns (is replaced by) e.Expr",
		Category:    "I/O Functions",
		Signature:   "<Print $0e.Expr>",
		ReturnValue: "None",
	},
	{
		Name: "Prout",
		Description: `Prints the expression e.Expr on the current output and returns the empty expression.
Functions that work with files require a file descriptor as an argument. A file descriptor is a macrodigit in the range 1-19; in some operations the descriptor 0 is allowed and refers to the terminal.`,
		Category:    "I/O Functions",
		Signature:   "<Prout $0e.Expr>",
		ReturnValue: "None",
	},
}

var Keywords = []Keyword{
	{
		Name:        "EXTERNAL",
		Value:       "$EXTERNAL",
		Description: "keyword $EXTERNAL",
	},
	{
		Name:        "EXTERN",
		Value:       "$EXTERN",
		Description: "keyword $EXTERN",
	},
	{
		Name:        "EXT",
		Value:       "$EXT",
		Description: "keyword $EXT",
	},
	{
		Name:        "ENTRY",
		Value:       "$ENTRY",
		Description: "keyword $ENTRY",
	},
}
