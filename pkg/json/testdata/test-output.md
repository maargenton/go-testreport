
# Test report


## github.com/maargenton/go-cli/pkg/cli

Coverage: 100%

- ✅ TestSetProcessEnv    
  - ✅ Given an environment as a list of strings    
    - ✅ when calling cmd.SetProcessEnv    
      - ✅ then all environment values are recorded  
- ✅ TestCommandRun    
  - ✅ Given a well defined command struct    
    - ✅ when calling run with valid arguments    
      - ✅ then the fields are set and the command is run    
    - ✅ when the command handler returns an error    
      - ✅ then the error is returned    
    - ✅ when calling run with invalid arguments    
      - ✅ then an error is returned    
    - ✅ when calling run with invalid environment value    
      - ✅ then an error is returned    
  - ✅ Given a command with invalid default    
    - ✅ when calling run with valid arguments    
      - ✅ then an error is returned    
  - ✅ Given a command with bad option    
    - ✅ when calling run with valid arguments    
      - ✅ then an error is returned  
- ✅ TestCommandUsage    
  - ✅ Given a well defined command struct    
    - ✅ when calling Usage()    
      - ✅ then a command description is returned    
  - ✅ Given an invalid command struct    
    - ✅ when calling Usage()    
      - ✅ then the error is printed in place of the usage  
- ✅ TestCustomUsage    
  - ✅ Given a command struct with custom usage    
    - ✅ when calling Usage()    
      - ✅ then the error is printed in place of the usage  
- ✅ TestCommandVersion    
  - ✅ Given a command struct with version handler    
    - ✅ when calling Usage()    
      - ✅ then a command description is returned  
- ✅ TestBashCompletionScript  
- ✅ TestFormatCompletionSuggestions    
  - ✅ Given a list of suggestions    
    - ✅ when calling FormatCompletionSuggestions()    
      - ✅ then all lines are at most the width    
    - ✅ when calling FormatCompletionSuggestions() with a single option    
      - ✅ then only the first word of the option is printed  
- ✅ TestDefaultCompletion    
  - ✅ Given the current directory structure    
    - ✅ when Calling DefaultCompletion() with an empty string    
      - ✅ then suggestions include the local files    
    - ✅ when Calling DefaultCompletion() with partial filename    
      - ✅ then suggestions include only the matching files    
    - ✅ when Calling DefaultCompletion() with partial unique folder name    
      - ✅ then suggestions include the files in that folder  
- ✅ TestFilepathCompletion    
  - ✅ Given a call to FilepathCompletion()    
    - ✅ when passing a pattern and an empty string    
      - ✅ then suggestions include all filenames matching the pattern    
    - ✅ when passing a pattern and a partial name    
      - ✅ then suggestions include only filenames matching both    
    - ✅ when passing a pattern and a non-matching partial name    
      - ✅ then the pattern is ignored and all matching files are returned  
- ✅ TestCommandRunCompletion    
  - ✅ Given a well defined command struct    
    - ✅ when calling Run() with completion request and partial option flag    
      - ✅ then the command is not run    
      - ✅ then the completion request error is returned    
      - ✅ then the suggestions contain the matching flag    
    - ✅ when calling Run() with a partial option argument    
      - ✅ then the command is not run    
      - ✅ then the completion request error is returned    
      - ✅ then the suggestions contain matching local filenames    
    - ✅ when calling Run() with nothing    
      - ✅ then the command is not run    
      - ✅ then the completion request error is returned    
      - ✅ then the suggestions include option flags    
      - ✅ then the suggestions include argument options    
  - ✅ Given a command with custom completion handler    
    - ✅ when calling Run() with nothing    
      - ✅ then the command is not run    
      - ✅ then the completion request error is returned    
      - ✅ then the suggestions include option flags    
      - ✅ then the suggestions include argument options    
    - ✅ when calling Run() with missing option argument    
      - ✅ then the suggestions contain matching local filenames  


## github.com/maargenton/go-cli/pkg/option

Coverage: 100%

- ✅ TestLineWrap    
  - ✅ Given a piece of text wider than allowed    
    - ✅ when calling lineWrap()    
      - ✅ then all lines are at most the width    
    - ✅ when calling lineWrap() with narrow width    
      - ✅ then lines longer than width have no space    
  - ✅ Given a piece of text with line breaks    
    - ✅ when calling lineWrap()    
      - ✅ then line breaks are preserved  
- ✅ TestParseOptsTag  
- ✅ TestParseOptsTagPartial  
- ✅ TestParseOptsTagInvalid  
- ✅ TestSplitSliceValues    
  - ✅ a,b,c    
  - ✅ a\,b,c    
  - ✅ a,b,\c    
  - ✅ a\\,b,c    
  - ✅ a\\:b;c  
- ✅ TestUnescapeField    
  - ✅ ab\c    
  - ✅ abc\\    
  - ✅ abc\    
  - ✅ abc\\def  
- ✅ TestScanTagFields    
  - ✅ short and long    
  - ✅ short and long with spaces    
  - ✅ with env default    
  - ✅ with env default and spaces    
  - ✅ with escape  
- ❌ TestGetCompletion    
  - ❌ Given an configured OptionSet    
    - ✅ when calling GetCompletion() with no arguments    
      - ✅ then no specific option is being completed    
      - ✅ then the first argument is being completed    
      - ✅ then available options are listed    
    - ✅ when calling GetCompletion() with flag expecting a value    
      - ✅ then only yhe option is returned    
    - ❌ when calling GetCompletion() with flag and value    
      - ❌ then remaining options are listed  
        ```
        completion_test.go:71:
            expected:       set(value.Option) == []string{ "--format <value>", "-v" }
            value:          []option.Description{
                            	option.Description{
                            		Option:      "--forfmat <value>",
                            		Description: "communcation format, e.g. 8N1",
                            	},
                            	option.Description{
                            		Option:      "-v",
                            		Description: "display additional information on startup",
                            	},
                            }
            $.Option:       { "--forfmat <value>", "-v" }
            extra values:   "--forfmat <value>"
            missing values: "--format <value>"
        ```  
      - ✅ then the first argument is being completed    
    - ❌ when calling GetCompletion() with short flag and value    
      - ❌ then remaining options are listed  
        ```
        completion_test.go:86:
            expected:       set(value.Option) == []string{ "--format <value>", "-v" }
            value:          []option.Description{
                            	option.Description{
                            		Option:      "--forfmat <value>",
                            		Description: "communcation format, e.g. 8N1",
                            	},
                            	option.Description{
                            		Option:      "-v",
                            		Description: "display additional information on startup",
                            	},
                            }
            $.Option:       { "--forfmat <value>", "-v" }
            extra values:   "--forfmat <value>"
            missing values: "--format <value>"
        ```  
      - ✅ then the first argument is being completed    
    - ✅ when calling GetCompletion() with exclusive flag    
      - ✅ then remaining options are listed    
    - ❌ when calling GetCompletion() with argument and flags    
      - ❌ then remaining options are listed  
        ```
        completion_test.go:113:
            expected:       set(value.Option) == []string{ "--format <value>", "-v", "--timestamp" }
            value:          []option.Description{
                            	option.Description{
                            		Option:      "--forfmat <value>",
                            		Description: "communcation format, e.g. 8N1",
                            	},
                            	option.Description{
                            		Option:      "--timestamp",
                            		Description: "prefix every line with elapsed time",
                            	},
                            	option.Description{
                            		Option:      "-v",
                            		Description: "display additional information on startup",
                            	},
                            }
            $.Option:       { "--forfmat <value>", "--timestamp", "-v" }
            extra values:   "--forfmat <value>"
            missing values: "--format <value>"
        ```  
      - ✅ then next argument is being completed  
- ✅ TestFormatOptionDescription    
  - ✅ Given a list of argument usage    
    - ✅ when calling FormatOptionDescription()    
      - ✅ then all descriptions are aligned    
  - ✅ Given a list of argument usage with long descriptions    
    - ✅ when calling FormatOptionDescription()    
      - ✅ then descriptions is wrapped and aligned accros multiple lines  
- ✅ TestFormatCompletion    
  - ✅ Given a list of completion suggestions    
    - ✅ when calling FormatCompletion()    
      - ✅ then suggestions are formatted one per line    
    - ✅ when calling FormatCompletion() with narrow width    
      - ✅ then long descriptions are truncated    
    - ✅ when calling FormatCompletion() with extremely narrow width    
      - ✅ then descriptions are dropped  
- ✅ TestOptionName    
  - ✅ Given a struct with opts tags    
    - ✅ when getting the name of a field with long flag    
      - ✅ the name includes the long flag    
    - ✅ when getting the name of a field with only short flag    
      - ✅ the name includes the short flag    
    - ✅ when getting the name of a positional argument field    
      - ✅ the name includes the value name    
      - ✅ the name includes the argument position    
    - ✅ when getting the name of an extra arguments field    
      - ✅ the name includes the value name    
  - ✅ Given another struct with opts tags    
    - ✅ when getting the name of an extra arguments field    
      - ✅ the name includes default name and ellipsis  
- ✅ TestOptionDescriptionUsage    
  - ✅ Given an Option{} with short and long    
    - ✅ when calling Usage()    
      - ✅ then it returns formatted usage    
  - ✅ Given an Option{} with short only    
    - ✅ when calling Usage()    
      - ✅ then it returns formatted usage    
  - ✅ Given an Option{} with long only    
    - ✅ when calling Usage()    
      - ✅ then it returns formatted usage    
  - ✅ Given an Option{} with named value    
    - ✅ when calling Usage()    
      - ✅ then it returns formatted usage    
  - ✅ Given an positional Option{} with nameonly    
    - ✅ when calling Usage()    
      - ✅ then it returns formatted usage  
- ✅ TestOptionDescriptionDescription    
  - ✅ Given an Option{} with no description    
    - ✅ when calling Usage()    
      - ✅ then it returns formatted description    
  - ✅ Given an Option{} with description    
    - ✅ when calling Usage()    
      - ✅ then it returns formatted description    
  - ✅ Given an Option{} with default and env    
    - ✅ when calling Usage()    
      - ✅ then it returns formatted description    
  - ✅ Given an Option{} with only default and env    
    - ✅ when calling Usage()    
      - ✅ then it returns formatted description  
- ✅ TestOption SetBool    
  - ✅ Given a struct with bool options    
    - ✅ when calling Set() on a bool field    
      - ✅ then the field is set to true    
    - ✅ when calling Set() on a bool pointer field    
      - ✅ then the field is set to point to a true value    
    - ✅ when calling Set() on a non-bool field    
      - ✅ then it panics  
- ✅ TestOption SetValue    
  - ✅ Given a struct with `opts` tags    
    - ✅ when setting the value of a scalar field    
      - ✅ with a valid value    
        - ✅ then the field is set accordingly    
      - ✅ with an invalid value    
        - ✅ then an error is returned and the value is not changed    
    - ✅ when setting the value of a pointer field    
      - ✅ with a valid value    
        - ✅ then the field is set accordingly    
      - ✅ with an invalid value    
        - ✅ then an error is returned and the value is not changed    
    - ✅ when setting the value of a slice field    
      - ✅ with a single valid value    
        - ✅ then the field is set accordingly    
      - ✅ with multiple delimited values    
        - ✅ then all values are recorded    
      - ✅ with empty value    
        - ✅ then all previous values are deleted    
      - ✅ with invalid value    
        - ✅ then all values are recorded  
- ✅ TestOptionSet  
- ✅ TestNewOptionSet NonStruct  
- ✅ TestNewOptionSet InvalidFieldType  
- ✅ TestNewOptionSet InvalidTag  
- ✅ TestNewOptionSet NestedInvalidTag  
- ✅ TestNewOptionSet UnexportedField  
- ✅ TestNewOptionSet WithEmptyOptsTag  
- ✅ TestNewOptionSet ArgN  
- ✅ TestNewOptionSet ArgN BadType  
- ✅ TestNewOptionSet Arg0  
- ✅ TestNewOptionSet ArgNonInt  
- ✅ TestNewOptionSet ArgGaps  
- ✅ TestNewOptionSet ArgDup  
- ✅ TestNewOptionSet Args  
- ✅ TestNewOptionSet ArgsDup  
- ✅ TestNewOptionSet ArgsNotSlice  
- ✅ TestOption    
  - ✅ d    
  - ✅ duration    
  - ✅ n    
  - ✅ name    
  - ✅ q  
- ✅ TestAddSpecialFlag    
  - ✅ Given an configured OptionSet    
    - ✅ when adding special flag with conflicting short    
      - ✅ then short flag is up-cased    
    - ✅ when adding special flag with conflicting short and up-cased short    
      - ✅ then short flag is dropped    
    - ✅ when adding special flag with conflicting long    
      - ✅ then the existing flag is preserved  
- ✅ TestApplyDefaults  
- ✅ TestApplyDefaultsError  
- ✅ TestApplyEnv  
- ✅ TestApplyEnvError  
- ✅ TestApplyArgsCombiningBoolFlags    
  - ✅ -ac    
  - ✅ -a -c    
  - ✅ -fabc    
  - ✅ -abfabc    
  - ✅ -afabc -b    
  - ✅ -af abc -b  
- ✅ TestApplyArgsLongFlags    
  - ✅ --aaa --ccc --file abc    
  - ✅ --aaa=false --ccc --file=abc  
- ✅ TestApplyArgs Delimiter    
  - ✅ Given a struct with args field    
    - ✅ when arguments contain `--`    
      - ✅ then remaining flags are not parsed    
    - ✅ when option value is `--`    
      - ✅ then remaining flags are parsed  
- ✅ TestApplyArgs Errors    
  - ✅ --aaa --ddd    
  - ✅ -abgc    
  - ✅ -d4p    
  - ✅ -d 4p    
  - ✅ --duration 4p    
  - ✅ --duration=4p    
  - ✅ a --duration  
- ✅ TestApplyArgs WithSpecialFlags    
  - ✅ --aaa --help    
  - ✅ --aaa -h    
  - ✅ --aaa --version    
  - ✅ --aaa -v  
- ✅ TestApplyArgs NonFlagArguments    
  - ✅ --aaa --duration 5m aaa bbb ccc ddd  
- ✅ TestApplyArgs NonFlagArguments Errors    
  - ✅ --aaa aaa    
  - ✅ --aaa aaa bbb ccc    
  - ✅ --aaa 1 2 ccc  


## github.com/maargenton/go-cli/pkg/value

Coverage: 100%

- ✅ TestParseBool  
- ✅ TestParseInt  
- ✅ TestParseInt8  
- ✅ TestParseInt16  
- ✅ TestParseInt32  
- ✅ TestParseInt64  
- ✅ TestParseUInt  
- ✅ TestParseUInt8  
- ✅ TestParseUInt16  
- ✅ TestParseUInt32  
- ✅ TestParseUInt64  
- ✅ TestParseFloat32  
- ✅ TestParseFloat64  
- ✅ TestParseString  
- ✅ TestRegisterParserPanicsWithInvalidArgument  
- ✅ TestRegisterParserPanicsWithMultipleParserForTheSameType  
- ✅ TestFlagValue  
- ✅ TestFlagValueError  
- ✅ TestTextUnmarshaller  
- ✅ TestUnparsable  
- ✅ TestTimeDuration  
- ✅ TestParseArgument  
