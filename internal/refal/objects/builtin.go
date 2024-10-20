package objects

// <Open s.Mode s.D e.File-name>
// opens file e.File-name and associates it with the file descriptor s.D . s.Mode is one of: 'w','W' (open for writing), 'r','R' (open for reading) or 'a','A' (open for adding). e.File-name may be empty; Open will then try to open file REFALdescr.DAT, where descr is the decimal representation of s.D . If the mode is reading and this file does not exist, an error occurs. If the mode is writing, this file is created.
// <Get s.D>
// where s.D is a file descriptor or, is similar to <Card> except that it gets its input from the file associated with s.D . If no file has been opened for reading with that file descriptor, Get will try to open file REFALdescr.DAT, where descr is the decimal representation of s.D . If it fails, an error occurs and the program is terminated. If s.D is 0, Get will read from the terminal.
// <Put s.D e.Expr>
// where s.D is a file descriptor or, writes e.Expr on the file associated with s.D and returns Expr (similar to Print). If no file has been opened for writing with that file descriptor, Put will open file REFALdescr.DAT, where descr is the decimal representation of s.D . If s.D is 0, Put prints out on the terminal. (Note that this output is not redirectable.)
// <Putout s.D e.Expr>
// returns empty (like Prout ). In all other respects Putout is identical to Put.
// <RemoveFile e.file_name> :: s.Boolean (e.Errors)
// If the function e.file_name can be successfully deleted, it is deleted, and s.Boolean is "True"; otherwise s.Boolean is "False", and the value of e.Errors is a description of the error.

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
		Signature:   "<Close s.ID>",
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
		Signature:   "<Print e.Expr>",
		ReturnValue: "None",
	},
	{
		Name: "Prout",
		Description: `Prints the expression e.Expr on the current output and returns the empty expression.
Functions that work with files require a file descriptor as an argument. A file descriptor is a macrodigit in the range 1-19; in some operations the descriptor 0 is allowed and refers to the terminal.`,
		Category:    "I/O Functions",
		Signature:   "<Prout e.Expr>",
		ReturnValue: "None",
	},
}
