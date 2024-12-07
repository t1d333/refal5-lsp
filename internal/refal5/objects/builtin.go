package objects

// TODO: add preview for functions, e.g. <Close s.ID>, but not <Close  ${1:s.ID}>
var BuiltInFunctions = map[string]RefalFunction{
	// I/O Functions
	"Card": {
		Name:        "Card",
		Description: "Returns (is replaced by) the next line from the input. Normally it is from the terminal, but input can be redirected as allowed by MS-DOS. The returned expression is a sequence of character-symbols (possibly empty). The End-Of-Line byte is not included. If the input is read from a file, the macrodigit 0 is returned when the end of file is reached (no more lines). This is used in programs as the indicator of end, since a macrodigit cannot result from input otherwise. When reading from the terminal, enter Control-Z to produce the same effect",
		Category:    "I/O Functions",
		Signature:   "<Card>",
		ReturnValue: "Object",
	},
	"Close": {
		Name:        "Close",
		Description: "s.ID is the identifying number of the file. The function closes the file. If such a file does not exist, no action is taken. The value of Close is always empty",
		Category:    "I/O Functions",
		Signature:   "<Close  ${1:s.ID}>",
	},
	"ExistFile": {
		Name:        "ExistFile",
		Description: "e.name is a file name (a string of characters). The function checks whether the file exists and returnes the corresponding identifier \"True\" or \"False\"",
		Category:    "I/O Functions",
		Signature:   "",
		ReturnValue: "Boolean",
	},
	"Print": {
		Name:        "Print",
		Description: "prints the expression e.Expr on the current output and returns (is replaced by) e.Expr",
		Category:    "I/O Functions",
		Signature:   "<Print ${1:e.Expr}>",
		ReturnValue: "None",
	},
	"Prout": {
		Name: "Prout",
		Description: `Prints the expression e.Expr on the current output and returns the empty expression.
Functions that work with files require a file descriptor as an argument. A file descriptor is a macrodigit in the range 1-19; in some operations the descriptor 0 is allowed and refers to the terminal.`,
		Category:    "I/O Functions",
		Signature:   "<Prout ${1:e.Expr}>",
		ReturnValue: "None",
	},
	"Open": {
		Name:        "Open",
		Description: `Opens file e.File-name and associates it with the file descriptor s.D . s.Mode is one of: 'w','W' (open for writing), 'r','R' (open for reading) or 'a','A' (open for adding). e.File-name may be empty; Open will then try to open file REFALdescr.DAT, where descr is the decimal representation of s.D . If the mode is reading and this file does not exist, an error occurs. If the mode is writing, this file is created`,
		Category:    "I/O Functions",
		Signature:   "<Open ${1:s.Mode} ${2:s.D} ${3:e.File-name}>",
		ReturnValue: "None",
	},
	"Get": {
		Name:        "Get",
		Description: `Where s.D is a file descriptor or, is similar to <Card> except that it gets its input from the file associated with s.D . If no file has been opened for reading with that file descriptor, Get will try to open file REFALdescr.DAT, where descr is the decimal representation of s.D . If it fails, an error occurs and the program is terminated. If s.D is 0, Get will read from the terminal`,
		Category:    "I/O Functions",
		Signature:   "<Get ${1:s.D}>",
		ReturnValue: "None",
	},
	"Put": {
		Name:        "Put",
		Description: `Where s.D is a file descriptor or, writes e.Expr on the file associated with s.D and returns Expr (similar to Print). If no file has been opened for writing with that file descriptor, Put will open file REFALdescr.DAT, where descr is the decimal representation of s.D . If s.D is 0, Put prints out on the terminal. (Note that this output is not redirectable.)`,
		Category:    "I/O Functions",
		Signature:   "<Put ${1:s.D} ${2:e.Expr}>",
		ReturnValue: "None",
	},
	"Putout": {
		Name:        "Putout",
		Description: `returns empty (like Prout ). In all other respects Putout is identical to Put`,
		Category:    "I/O Functions",
		Signature:   "<Putout ${1:s.D} ${2:e.Expr}>",
		ReturnValue: "None",
	},
	"RemoveFile": {
		Name:        "RemoveFile",
		Description: `If the function e.file_name can be successfully deleted, it is deleted, and s.Boolean is "True"; otherwise s.Boolean is "False", and the value of e.Errors is a description of the error.`,
		Category:    "I/O Functions",
		Signature:   "<RemoveFile ${1:e.file_name}>",
		ReturnValue: "None",
	},

	// Arithmetic Functions
	"Add": {
		Name:        "Add",
		Description: `Returns the sum of the operands`,
		Category:    "Arithmetic Functions",
		Signature:   "<Add ${1:e.N1} ${2:e.N2}>",
		ReturnValue: "None",
	},

	"+": {
		Name:        "+",
		Description: `Returns the sum of the operands`,
		Category:    "Arithmetic Functions",
		Signature:   "<+ ${1:e.N1} ${2:e.N2}>",
		ReturnValue: "None",
	},
	"Sub": {
		Name:        "Sub",
		Description: `Subtracts e.N2 from e.N1 and returns the difference`,
		Category:    "Arithmetic Functions",
		Signature:   "<Sub ${1:e.N1} ${2:e.N2}>",
		ReturnValue: "None",
	},

	"-": {
		Name:        "-",
		Description: `Subtracts e.N2 from e.N1 and returns the difference`,
		Category:    "Arithmetic Functions",
		Signature:   "<- ${1:e.N1} ${2:e.N2}>",
		ReturnValue: "None",
	},

	"Mul": {
		Name:        "Mul",
		Description: `Returns the product of the operands`,
		Category:    "Arithmetic Functions",
		Signature:   "<Mul ${1:e.N1} ${2:e.N2}>",
		ReturnValue: "None",
	},

	"*": {
		Name:        "*",
		Description: `Returns the product of the operands`,
		Category:    "Arithmetic Functions",
		Signature:   "<Mul ${1:e.N1} ${2:e.N2}>",
		ReturnValue: "None",
	},

	"Div": {
		Name:        "Div",
		Description: `Returns the whole quotient of e.N1 and e.N2 ; the remainder is ignored. With this and the other two division functions, if e.N2 is 0, an error occurs`,
		Category:    "Arithmetic Functions",
		Signature:   "<Div ${1:e.N1} ${2:e.N2}>",
		ReturnValue: "None",
	},

	"/": {
		Name:        "/",
		Description: `Returns the whole quotient of e.N1 and e.N2 ; the remainder is ignored. With this and the other two division functions, if e.N2 is 0, an error occurs`,
		Category:    "Arithmetic Functions",
		Signature:   "</ ${1:e.N1} ${2:e.N2}>",
		ReturnValue: "None",
	},
	"Divmod": {
		Name:        "Divmod",
		Description: `Expects whole arguments and returns (e.Quotient) e.Remainder  The remainder is of the sign of e.N1 .`,
		Category:    "Arithmetic Functions",
		Signature:   "<Divmod ${1:e.N1} ${2:e.N2}>",
		ReturnValue: "None",
	},
	"Mod": {
		Name:        "Mod",
		Description: `expects whole arguments and returns the remainder of the division of e.N1 by e.N2 .`,
		Category:    "Arithmetic Functions",
		Signature:   "<Mod ${1:e.N1} ${2:e.N2}>",
		ReturnValue: "None",
	},
	"Compare": {
		Name:        "Compare",
		Description: `compares the two numbers and returns '-' if e.N1 is less than e.N2 , '+' if it is greater, and '0' if they are equal.`,
		Category:    "Arithmetic Functions",
		Signature:   "<Compare ${1:e.N1} ${2:e.N2}>",
		ReturnValue: "None",
	},

	// Stack Operations

	"Br": {
		Name:        "Br",
		Description: `Buries the expression e.Name under the name e.Name . The name must not contain '=' on the top level of structure`,
		Category:    "Stack Operations",
		Signature:   "<Br ${1:e.Name} '=' ${2:e.Expr}>",
		ReturnValue: "None",
	},

	"Dg": {
		Name:        "Dg",
		Description: `digs the expression buried under the name e.Name , i.e., returns the last expression buried under this name and removes it from the stack. If there is no expression buried under e.Name , Dg returns the empty expression`,
		Category:    "Stack Operations",
		Signature:   "<Dg ${1:e.Name}>",
		ReturnValue: "None",
	},

	"Cp": {
		Name:        "Cp",
		Description: `Works as Dg but does not remove the expression from the stack`,
		Category:    "Stack Operations",
		Signature:   "<Cp ${1:e.Name}>",
		ReturnValue: "None",
	},

	"Rp": {
		Name:        "Rp",
		Description: `replaces the expression buried under the name e.Name by e.Expr`,
		Category:    "Stack Operations",
		Signature:   "<Rp ${1:e.Name} '=' ${2:e.Expr}>",
		ReturnValue: "None",
	},

	"Dgall": {
		Name:        "Dgall",
		Description: `digs out the whole stack. The stack is a string of terms of the form (e.Name'=' e.Value) Each time Br is activated, such a term is added on the left side. Dg takes the leftmost term away, while Cp copies it, and Rp changes it.`,
		Category:    "Stack Operations",
		Signature:   "<Dgall>",
		ReturnValue: "None",
	},

	// Symbol and String Manipulation
	"Type": {
		Name:        "Type",
		Description: `Returns s.Type s.Subtype e.Expr , where e.Expr is unchanged and s.Type and s.Subtype depend on the type of the first element of e.Expr`,
		Category:    "Symbol and String Manipulation",
		Signature:   "<Type ${1:e.Expr}>",
		ReturnValue: "None",
	},

	"Numb": {
		Name:        "Numb",
		Description: `Returns the macrodigit represented by e.Digit-string`,
		Category:    "Symbol and String Manipulation",
		Signature:   "<Numb ${1:e.Digit-string}>",
		ReturnValue: "None",
	},

	"Symb": {
		Name:        "Symb",
		Description: `Is the inverse of Numb. It returns the string of decimal digits representing s.Macrodigit`,
		Category:    "Symbol and String Manipulation",
		Signature:   "<Symb ${1:s.Macrodigit}>",
		ReturnValue: "None",
	},

	"Implode": {
		Name:        "Implode",
		Description: `Takes the initial alphanumeric characters of e.Expr and creates an identifier (symbolic name) from them. The initial string in e.Expr must begin with a letter and terminate with a non-alphanumeric character, parenthesis, or the end of the expression. Underscore and dash are also permitted. Implode returns the identifier followed by the unimploded portion of e.Expr . If the first character is not a letter, Implode returns the macrodigit 0 followed by the argument`,
		Category:    "Symbol and String Manipulation",
		Signature:   "<Implode ${1:e.Expr}>",
		ReturnValue: "None",
	},

	"Explode": {
		Name:        "Explode",
		Description: `Returns the string of character-symbols which make up s.Compound_Symbol`,
		Category:    "Symbol and String Manipulation",
		Signature:   "<Explode ${1:s.Compound_Symbol}>",
		ReturnValue: "None",
	},

	"Chr": {
		Name:        "Chr",
		Description: `Replaces every macrodigit in e.Expr by the character-symbol with the same ASCII code modulo 256`,
		Category:    "Symbol and String Manipulation",
		Signature:   "<Chr ${1:e.Expr}>",
		ReturnValue: "None",
	},

	"Ord": {
		Name:        "Ord",
		Description: `Is the inverse of Char. It returns the expression in which all characters are replaced by macrodigits equal to their ASCII codes`,
		Category:    "Symbol and String Manipulation",
		Signature:   "<Ord ${1:e.Expr}>",
		ReturnValue: "None",
	},

	"First": {
		Name:        "First",
		Description: `Where s.N is a macrodigit, breaks up e.Expr into two parts — e.1 and e.2 , and returns (e.1)e.2 . If the original expression e.Expr has at least s.N terms (on the top level of structure), then the first s.N terms go to e.1 , and the rest to e.2 . Otherwise, e.1 is e.Expr and e.2 is empty`,
		Category:    "Symbol and String Manipulation",
		Signature:   "<First ${1:s.N} ${2:e.Expr}>",
		ReturnValue: "None",
	},

	"Last": {
		Name:        "Last",
		Description: `Is similar to First but it is e.2 that has s.N terms`,
		Category:    "Symbol and String Manipulation",
		Signature:   "<Last ${1:s.N} ${2:e.Expr}>",
		ReturnValue: "None",
	},

	"Lenw": {
		Name:        "Lenw",
		Description: `Returns the length of e.Expr in terms followed by e.Expr`,
		Category:    "Symbol and String Manipulation",
		Signature:   "<Lenw ${1:e.Expr}>",
		ReturnValue: "None",
	},

	"Lower": {
		Name:        "Lower",
		Description: `Returns the original expression e.Expr in which all capital letters are replaced by lower case letters`,
		Category:    "Symbol and String Manipulation",
		Signature:   "<Lower ${1:e.Expr}>",
		ReturnValue: "None",
	},

	"Upper": {
		Name:        "Upper",
		Description: `Is similar to Lower. All lower case letters are capitalized`,
		Category:    "Symbol and String Manipulation",
		Signature:   "<Upper ${1:e.Expr}>",
		ReturnValue: "None",
	},

	// System Functions
	"GetCurrentDirectory": {
		Name:        "GetCurrentDirectory",
		Description: `This function returns the name of the current directory where the Refal-5 interpreter is running`,
		Category:    "System Functions",
		Signature:   "<GetCurrentDirectory>",
		ReturnValue: "None",
	},

	"GetEnv": {
		Name:        "GetEnv",
		Description: `Takes characters and interprets them as a name of the environment variable of the current operating system. The string must not exceed 120 characters. If the variable is undefined, then the function returns the empty expression, otherwise it returns a value of the variable as a sequence of characters. If the argument is not a sequence of characters, then "Recognition impossible" occures`,
		Category:    "System Functions",
		Signature:   "<GetEnv ${1:e.characters}>",
		ReturnValue: "None",
	},

	"Exit": {
		Name:        "Exit",
		Description: `Takes a macrodigit or two terms ( '-' s.macrodigit ) and interrupts the current Refal-process. The process returns the argument to the current operating system as a return-code. If the argument is something else, then "Recognition impossible" occurs`,
		Category:    "System Functions",
		Signature:   "<Exit ${1:e.return-code}>",
		ReturnValue: "None",
	},

	"Step": {
		Name:        "Step",
		Description: `Returns the sequential number of the current step as a macrodigit`,
		Category:    "System Functions",
		Signature:   "<Step>",
		ReturnValue: "None",
	},

	"System": {
		Name:        "System",
		Description: `Takes e.characters and tries to interpret them as a shell's command of a current operating system. The string must not exceed 255 characters. If the value is a negative integer, then the function returns the string '-' followed by macrodigit, otherwise it returns a macrodigit. If the argument is not a sequence of characters, then "Recognition impossible" occurs`,
		Category:    "System Functions",
		Signature:   "<System ${1:e.characters}>",
		ReturnValue: "None",
	},

	"Time": {
		Name:        "Time",
		Description: `Returns a string containing the current system time`,
		Category:    "System Functions",
		Signature:   "<Time>",
		ReturnValue: "None",
	},

	"TimeElapsed": {
		Name:        "TimeElapsed",
		Description: `Takes empty expression or macrodigit zero. It returns elapsead system time. The time is elapsed time from a last call to the same function with zero as its argument. If there was no such call then the time is elapsed time from start of loading of RSL-modules. The time is given in seconds. If the argument is something else, then "Recognition impossible" occurs`,
		Category:    "System Functions",
		Signature:   "<TimeElapsed ${1:e.init}>",
		ReturnValue: "None",
	},

	"Arg": {
		Name:        "Arg",
		Description: `Where s.N is a macrodigit, returns the command line argument which has the sequential number s.N . Command line arguments must not start with '-' (in order not to be confused with flags)`,
		Category:    "System Functions",
		Signature:   "<Arg ${1:s.N}>",
		ReturnValue: "None",
	},

	"Up": {
		Name:        "Up",
		Description: `Upgrades e.Expr (demetacodes it). See documentation for restrictions on e.Expr`,
		Category:    "System Functions",
		Signature:   "<Up ${1:e.Expr}>",
		ReturnValue: "None",
	},

	"Dn": {
		Name:        "Dn",
		Description: `<Dn e.Expr> downgrades e.Expr (metacodes it)`,
		Category:    "System Functions",
		Signature:   "<Dn ${1:e.Expr}>",
		ReturnValue: "None",
	},

	"Mu": {
		Name:        "Mu",
		Description: `finds the function whose name is s.F-name or <Implode e.String> (when given in the form of a string of characters) and applies it to the expression e.Expr , i.e., is replaced by <s.F-name e.Expr>. If no such function is visible from the call of Mu, an error occurs.`,
		Category:    "System Functions",
		Signature:   "<Mu ${1:s.F-name} ${2:e.Expr}>",
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
