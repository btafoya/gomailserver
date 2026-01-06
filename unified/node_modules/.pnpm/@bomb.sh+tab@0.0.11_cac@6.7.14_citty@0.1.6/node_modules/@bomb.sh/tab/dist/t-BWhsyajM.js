import e from"node:assert";function t(e,t){return`#compdef ${e}
compdef _${e} ${e}

# zsh completion for ${e} -*- shell-script -*-

__${e}_debug() {
    local file="$BASH_COMP_DEBUG_FILE"
    if [[ -n \${file} ]]; then
        echo "$*" >> "\${file}"
    fi
}

_${e}() {
    local shellCompDirectiveError=${a.ShellCompDirectiveError}
    local shellCompDirectiveNoSpace=${a.ShellCompDirectiveNoSpace}
    local shellCompDirectiveNoFileComp=${a.ShellCompDirectiveNoFileComp}
    local shellCompDirectiveFilterFileExt=${a.ShellCompDirectiveFilterFileExt}
    local shellCompDirectiveFilterDirs=${a.ShellCompDirectiveFilterDirs}
    local shellCompDirectiveKeepOrder=${a.ShellCompDirectiveKeepOrder}

    local lastParam lastChar flagPrefix requestComp out directive comp lastComp noSpace keepOrder
    local -a completions

    __${e}_debug "\\n========= starting completion logic =========="
    __${e}_debug "CURRENT: \${CURRENT}, words[*]: \${words[*]}"

    # The user could have moved the cursor backwards on the command-line.
    # We need to trigger completion from the $CURRENT location, so we need
    # to truncate the command-line ($words) up to the $CURRENT location.
    # (We cannot use $CURSOR as its value does not work when a command is an alias.)
    words=( "\${=words[1,CURRENT]}" )
    __${e}_debug "Truncated words[*]: \${words[*]},"

    lastParam=\${words[-1]}
    lastChar=\${lastParam[-1]}
    __${e}_debug "lastParam: \${lastParam}, lastChar: \${lastChar}"

    # For zsh, when completing a flag with an = (e.g., ${e} -n=<TAB>)
    # completions must be prefixed with the flag
    setopt local_options BASH_REMATCH
    if [[ "\${lastParam}" =~ '-.*=' ]]; then
        # We are dealing with a flag with an =
        flagPrefix="-P \${BASH_REMATCH}"
    fi

    # Prepare the command to obtain completions, ensuring arguments are quoted for eval
    local -a args_to_quote=("\${(@)words[2,-1]}")
    if [ "\${lastChar}" = "" ]; then
        # If the last parameter is complete (there is a space following it)
        # We add an extra empty parameter so we can indicate this to the go completion code.
        __${e}_debug "Adding extra empty parameter"
        args_to_quote+=("")
    fi

    # Use Zsh's (q) flag to quote each argument safely for eval
    local quoted_args=("\${(@q)args_to_quote}")

    # Join the main command and the quoted arguments into a single string for eval
    requestComp="${t} complete -- \${quoted_args[*]}"

    __${e}_debug "About to call: eval \${requestComp}"

    # Use eval to handle any environment variables and such
    out=$(eval \${requestComp} 2>/dev/null)
    __${e}_debug "completion output: \${out}"

    # Extract the directive integer following a : from the last line
    local lastLine
    while IFS='\n' read -r line; do
        lastLine=\${line}
    done < <(printf "%s\n" "\${out[@]}")
    __${e}_debug "last line: \${lastLine}"

    if [ "\${lastLine[1]}" = : ]; then
        directive=\${lastLine[2,-1]}
        # Remove the directive including the : and the newline
        local suffix
        (( suffix=\${#lastLine}+2))
        out=\${out[1,-$suffix]}
    else
        # There is no directive specified.  Leave $out as is.
        __${e}_debug "No directive found.  Setting to default"
        directive=0
    fi

    __${e}_debug "directive: \${directive}"
    __${e}_debug "completions: \${out}"
    __${e}_debug "flagPrefix: \${flagPrefix}"

    if [ $((directive & shellCompDirectiveError)) -ne 0 ]; then
        __${e}_debug "Completion received error. Ignoring completions."
        return
    fi

    local activeHelpMarker="%"
    local endIndex=\${#activeHelpMarker}
    local startIndex=$((\${#activeHelpMarker}+1))
    local hasActiveHelp=0
    while IFS='\n' read -r comp; do
        # Check if this is an activeHelp statement (i.e., prefixed with $activeHelpMarker)
        if [ "\${comp[1,$endIndex]}" = "$activeHelpMarker" ];then
            __${e}_debug "ActiveHelp found: $comp"
            comp="\${comp[$startIndex,-1]}"
            if [ -n "$comp" ]; then
                compadd -x "\${comp}"
                __${e}_debug "ActiveHelp will need delimiter"
                hasActiveHelp=1
            fi
            continue
        fi

        if [ -n "$comp" ]; then
            # If requested, completions are returned with a description.
            # The description is preceded by a TAB character.
            # For zsh's _describe, we need to use a : instead of a TAB.
            # We first need to escape any : as part of the completion itself.
            comp=\${comp//:/\\:}

            local tab="$(printf '\\t')"
            comp=\${comp//$tab/:}

            __${e}_debug "Adding completion: \${comp}"
            completions+=\${comp}
            lastComp=$comp
        fi
    done < <(printf "%s\n" "\${out[@]}")

    # Add a delimiter after the activeHelp statements, but only if:
    # - there are completions following the activeHelp statements, or
    # - file completion will be performed (so there will be choices after the activeHelp)
    if [ $hasActiveHelp -eq 1 ]; then
        if [ \${#completions} -ne 0 ] || [ $((directive & shellCompDirectiveNoFileComp)) -eq 0 ]; then
            __${e}_debug "Adding activeHelp delimiter"
            compadd -x "--"
            hasActiveHelp=0
        fi
    fi

    if [ $((directive & shellCompDirectiveNoSpace)) -ne 0 ]; then
        __${e}_debug "Activating nospace."
        noSpace="-S ''"
    fi

    if [ $((directive & shellCompDirectiveKeepOrder)) -ne 0 ]; then
        __${e}_debug "Activating keep order."
        keepOrder="-V"
    fi

    if [ $((directive & shellCompDirectiveFilterFileExt)) -ne 0 ]; then
        # File extension filtering
        local filteringCmd
        filteringCmd='_files'
        for filter in \${completions[@]}; do
            if [ \${filter[1]} != '*' ]; then
                # zsh requires a glob pattern to do file filtering
                filter="\\*.$filter"
            fi
            filteringCmd+=" -g $filter"
        done
        filteringCmd+=" \${flagPrefix}"

        __${e}_debug "File filtering command: $filteringCmd"
        _arguments '*:filename:'"$filteringCmd"
    elif [ $((directive & shellCompDirectiveFilterDirs)) -ne 0 ]; then
        # File completion for directories only
        local subdir
        subdir="\${completions[1]}"
        if [ -n "$subdir" ]; then
            __${e}_debug "Listing directories in $subdir"
            pushd "\${subdir}" >/dev/null 2>&1
        else
            __${e}_debug "Listing directories in ."
        fi

        local result
        _arguments '*:dirname:_files -/'" \${flagPrefix}"
        result=$?
        if [ -n "$subdir" ]; then
            popd >/dev/null 2>&1
        fi
        return $result
    else
        __${e}_debug "Calling _describe"
        if eval _describe $keepOrder "completions" completions -Q \${flagPrefix} \${noSpace}; then
            __${e}_debug "_describe found some completions"

            # Return the success of having called _describe
            return 0
        else
            __${e}_debug "_describe did not find completions."
            __${e}_debug "Checking if we should do file completion."
            if [ $((directive & shellCompDirectiveNoFileComp)) -ne 0 ]; then
                __${e}_debug "deactivating file completion"

                # Return 0 to indicate completion is finished and prevent zsh from
                # trying other completion algorithms (which could cause hanging).
                # We use NoFileComp directive to explicitly disable file completion.
                return 0
            else
                # Perform file completion
                __${e}_debug "Activating file completion"

                # We must return the result of this command, so it must be the
                # last command, or else we must store its result to return it.
                _arguments '*:filename:_files'" \${flagPrefix}"
            fi
        fi
    fi
}

# don't run the completion function when being sourced or eval-ed
if [ "\${funcstack[1]}" = "_${e}" ]; then
    _${e}
fi
`}function n(e,t){let n=e.replace(/[-:]/g,`_`);return`# bash completion for ${e}

# Define shell completion directives
readonly ShellCompDirectiveError=${a.ShellCompDirectiveError}
readonly ShellCompDirectiveNoSpace=${a.ShellCompDirectiveNoSpace}
readonly ShellCompDirectiveNoFileComp=${a.ShellCompDirectiveNoFileComp}
readonly ShellCompDirectiveFilterFileExt=${a.ShellCompDirectiveFilterFileExt}
readonly ShellCompDirectiveFilterDirs=${a.ShellCompDirectiveFilterDirs}
readonly ShellCompDirectiveKeepOrder=${a.ShellCompDirectiveKeepOrder}

# Function to debug completion
__${n}_debug() {
    if [[ -n \${BASH_COMP_DEBUG_FILE:-} ]]; then
        echo "$*" >> "\${BASH_COMP_DEBUG_FILE}"
    fi
}

# Function to handle completions
__${n}_complete() {
    local cur prev words cword
    _get_comp_words_by_ref -n "=:" cur prev words cword

    local requestComp out directive
    
    # Build the command to get completions
    requestComp="${t} complete -- \${words[@]:1}"
    
    # Add an empty parameter if the last parameter is complete
    if [[ -z "$cur" ]]; then
        requestComp="$requestComp ''"
    fi
    
    # Get completions from the program
    out=$(eval "$requestComp" 2>/dev/null)
    
    # Extract directive if present
    directive=0
    if [[ "$out" == *:* ]]; then
        directive=\${out##*:}
        out=\${out%:*}
    fi
    
    # Process completions based on directive
    if [[ $((directive & $ShellCompDirectiveError)) -ne 0 ]]; then
        # Error, no completion
        return
    fi
    
    # Apply directives
    if [[ $((directive & $ShellCompDirectiveNoSpace)) -ne 0 ]]; then
        compopt -o nospace
    fi
    if [[ $((directive & $ShellCompDirectiveKeepOrder)) -ne 0 ]]; then
        compopt -o nosort
    fi
    if [[ $((directive & $ShellCompDirectiveNoFileComp)) -ne 0 ]]; then
        compopt +o default
    fi
    
    # Handle file extension filtering
    if [[ $((directive & $ShellCompDirectiveFilterFileExt)) -ne 0 ]]; then
        local filter=""
        for ext in $out; do
            filter="$filter|$ext"
        done
        filter="\\.($filter)"
        compopt -o filenames
        COMPREPLY=( $(compgen -f -X "!$filter" -- "$cur") )
        return
    fi
    
    # Handle directory filtering
    if [[ $((directive & $ShellCompDirectiveFilterDirs)) -ne 0 ]]; then
        compopt -o dirnames
        COMPREPLY=( $(compgen -d -- "$cur") )
        return
    fi
    
    # Process completions
    local IFS=$'\\n'
    local tab=$(printf '\\t')
    
    # Parse completions with descriptions
    local completions=()
    while read -r comp; do
        if [[ "$comp" == *$tab* ]]; then
            # Split completion and description
            local value=\${comp%%$tab*}
            local desc=\${comp#*$tab}
            completions+=("$value")
        else
            completions+=("$comp")
        fi
    done <<< "$out"
    
    # Return completions
    COMPREPLY=( $(compgen -W "\${completions[*]}" -- "$cur") )
}

# Register completion function
complete -F __${n}_complete ${e}
`}function r(e,t){let n=e.replace(/[-:]/g,`_`);return`# fish completion for ${e} -*- shell-script -*-

# Define shell completion directives
set -l ShellCompDirectiveError ${a.ShellCompDirectiveError}
set -l ShellCompDirectiveNoSpace ${a.ShellCompDirectiveNoSpace}
set -l ShellCompDirectiveNoFileComp ${a.ShellCompDirectiveNoFileComp}
set -l ShellCompDirectiveFilterFileExt ${a.ShellCompDirectiveFilterFileExt}
set -l ShellCompDirectiveFilterDirs ${a.ShellCompDirectiveFilterDirs}
set -l ShellCompDirectiveKeepOrder ${a.ShellCompDirectiveKeepOrder}

function __${n}_debug
    set -l file "$BASH_COMP_DEBUG_FILE"
    if test -n "$file"
        echo "$argv" >> $file
    end
end

function __${n}_perform_completion
    __${n}_debug "Starting __${n}_perform_completion"

    # Extract all args except the completion flag
    set -l args (string match -v -- "--completion=" (commandline -opc))
    
    # Extract the current token being completed
    set -l current_token (commandline -ct)
    
    # Check if current token starts with a dash
    set -l flag_prefix ""
    if string match -q -- "-*" $current_token
        set flag_prefix "--flag="
    end
    
    __${n}_debug "Current token: $current_token"
    __${n}_debug "All args: $args"

    # Call the completion program and get the results
    set -l requestComp "${t} complete -- $args"
    __${n}_debug "Calling $requestComp"
    set -l results (eval $requestComp 2> /dev/null)
    
    # Some programs may output extra empty lines after the directive.
    # Let's ignore them or else it will break completion.
    # Ref: https://github.com/spf13/cobra/issues/1279
    for line in $results[-1..1]
        if test (string sub -s 1 -l 1 -- $line) = ":"
            # The directive
            set -l directive (string sub -s 2 -- $line)
            set -l directive_num (math $directive)
            break
        end
    end
    
    # No directive specified, use default
    if not set -q directive_num
        set directive_num 0
    end
    
    __${n}_debug "Directive: $directive_num"

    # Process completions based on directive
    if test $directive_num -eq $ShellCompDirectiveError
        # Error code. No completion.
        __${n}_debug "Received error directive: aborting."
        return 1
    end

    # Filter out the directive (last line)
    if test (count $results) -gt 0 -a (string sub -s 1 -l 1 -- $results[-1]) = ":"
        set results $results[1..-2]
    end

    # No completions, let fish handle file completions unless forbidden
    if test (count $results) -eq 0
        if test $directive_num -ne $ShellCompDirectiveNoFileComp
            __${n}_debug "No completions, performing file completion"
            return 1
        end
        __${n}_debug "No completions, but file completion forbidden"
        return 0
    end

    # Filter file extensions
    if test $directive_num -eq $ShellCompDirectiveFilterFileExt
        __${n}_debug "File extension filtering"
        set -l file_extensions
        for item in $results
            if test -n "$item" -a (string sub -s 1 -l 1 -- $item) != "-"
                set -a file_extensions "*$item"
            end
        end
        __${n}_debug "File extensions: $file_extensions"
        
        # Use the file extensions as completions
        set -l completions
        for ext in $file_extensions
            # Get all files matching the extension
            set -a completions (string replace -r '^.*/' '' -- $ext)
        end
        
        for item in $completions
            echo -e "$item\t"
        end
        return 0
    end

    # Filter directories
    if test $directive_num -eq $ShellCompDirectiveFilterDirs
        __${n}_debug "Directory filtering"
        set -l dirs
        for item in $results
            if test -d "$item"
                set -a dirs "$item/"
            end
        end
        
        for item in $dirs
            echo -e "$item\t"
        end
        return 0
    end

    # Process remaining completions
    for item in $results
        if test -n "$item"
            # Check if the item has a description
            if string match -q "*\t*" -- "$item"
                set -l completion_parts (string split \t -- "$item")
                set -l comp $completion_parts[1]
                set -l desc $completion_parts[2]
                
                # Add the completion and description
                echo -e "$comp\t$desc"
            else
                # Add just the completion
                echo -e "$item\t"
            end
        end
    end
    
    # If directive contains NoSpace, tell fish not to add a space after completion
    if test (math "$directive_num & $ShellCompDirectiveNoSpace") -ne 0
        return 2
    end
    
    return 0
end

# Set up the completion for the ${e} command
complete -c ${e} -f -a "(eval __${n}_perform_completion)"
`}function i(e,t){let n=e.replace(/[-:]/g,`_`);return`# powershell completion for ${e} -*- shell-script -*-

  [Console]::OutputEncoding = [System.Text.Encoding]::UTF8
    function __${e}_debug {
        if ($env:BASH_COMP_DEBUG_FILE) {
            "$args" | Out-File -Append -FilePath "$env:BASH_COMP_DEBUG_FILE"
        }
    }

    filter __${e}_escapeStringWithSpecialChars {
        $_ -replace '\\s|#|@|\\$|;|,|''|\\{|\\}|\\(|\\)|"|\\||<|>|&','\`$&'
    }

[scriptblock]$__${n}CompleterBlock = {
    param(
            $WordToComplete,
            $CommandAst,
            $CursorPosition
        )

    # Get the current command line and convert into a string
    $Command = $CommandAst.CommandElements
    $Command = "$Command"

    __${e}_debug ""
    __${e}_debug "========= starting completion logic =========="
    __${e}_debug "WordToComplete: $WordToComplete Command: $Command CursorPosition: $CursorPosition"

    # The user could have moved the cursor backwards on the command-line.
    # We need to trigger completion from the $CursorPosition location, so we need
    # to truncate the command-line ($Command) up to the $CursorPosition location.
    # Make sure the $Command is longer then the $CursorPosition before we truncate.
    # This happens because the $Command does not include the last space.
    if ($Command.Length -gt $CursorPosition) {
        $Command = $Command.Substring(0, $CursorPosition)
    }
    __${e}_debug "Truncated command: $Command"

    $ShellCompDirectiveError=${a.ShellCompDirectiveError}
    $ShellCompDirectiveNoSpace=${a.ShellCompDirectiveNoSpace}
    $ShellCompDirectiveNoFileComp=${a.ShellCompDirectiveNoFileComp}
    $ShellCompDirectiveFilterFileExt=${a.ShellCompDirectiveFilterFileExt}
    $ShellCompDirectiveFilterDirs=${a.ShellCompDirectiveFilterDirs}
    $ShellCompDirectiveKeepOrder=${a.ShellCompDirectiveKeepOrder}

    # Prepare the command to request completions for the program.
    # Split the command at the first space to separate the program and arguments.
    $Program, $Arguments = $Command.Split(" ", 2)

    $QuotedArgs = ($Arguments -split ' ' | ForEach-Object { "'" + ($_ -replace "'", "''") + "'" }) -join ' '
    __${e}_debug "QuotedArgs: $QuotedArgs"

    $RequestComp = "& ${t} complete '--' $QuotedArgs"
    __${e}_debug "RequestComp: $RequestComp"

    # we cannot use $WordToComplete because it
    # has the wrong values if the cursor was moved
    # so use the last argument
    if ($WordToComplete -ne "" ) {
        $WordToComplete = $Arguments.Split(" ")[-1]
    }
    __${e}_debug "New WordToComplete: $WordToComplete"


    # Check for flag with equal sign
    $IsEqualFlag = ($WordToComplete -Like "--*=*" )
    if ( $IsEqualFlag ) {
        __${e}_debug "Completing equal sign flag"
        # Remove the flag part
        $Flag, $WordToComplete = $WordToComplete.Split("=", 2)
    }

    if ( $WordToComplete -eq "" -And ( -Not $IsEqualFlag )) {
        # If the last parameter is complete (there is a space following it)
        # We add an extra empty parameter so we can indicate this to the go method.
        __${e}_debug "Adding extra empty parameter"
        # PowerShell 7.2+ changed the way how the arguments are passed to executables,
        # so for pre-7.2 or when Legacy argument passing is enabled we need to use
        if ($PSVersionTable.PsVersion -lt [version]'7.2.0' -or
            ($PSVersionTable.PsVersion -lt [version]'7.3.0' -and -not [ExperimentalFeature]::IsEnabled("PSNativeCommandArgumentPassing")) -or
            (($PSVersionTable.PsVersion -ge [version]'7.3.0' -or [ExperimentalFeature]::IsEnabled("PSNativeCommandArgumentPassing")) -and
              $PSNativeCommandArgumentPassing -eq 'Legacy')) {
             $RequestComp="$RequestComp" + ' \`"\`"'
        } else {
             $RequestComp = "$RequestComp" + ' ""'
        }
    }

    __${e}_debug "Calling $RequestComp"
    # First disable ActiveHelp which is not supported for Powershell
    $env:ActiveHelp = 0

    # call the command store the output in $out and redirect stderr and stdout to null
    # $Out is an array contains each line per element
    Invoke-Expression -OutVariable out "$RequestComp" 2>&1 | Out-Null

    # get directive from last line
    [int]$Directive = $Out[-1].TrimStart(':')
    if ($Directive -eq "") {
        # There is no directive specified
        $Directive = 0
    }
    __${e}_debug "The completion directive is: $Directive"

    # remove directive (last element) from out
    $Out = $Out | Where-Object { $_ -ne $Out[-1] }
    __${e}_debug "The completions are: $Out"

    if (($Directive -band $ShellCompDirectiveError) -ne 0 ) {
        # Error code.  No completion.
        __${e}_debug "Received error from custom completion go code"
        return
    }

    $Longest = 0
    [Array]$Values = $Out | ForEach-Object {
        # Split the output in name and description
        $Name, $Description = $_.Split("\`t", 2)
        __${e}_debug "Name: $Name Description: $Description"

        # Look for the longest completion so that we can format things nicely
        if ($Longest -lt $Name.Length) {
            $Longest = $Name.Length
        }

        # Set the description to a one space string if there is none set.
        # This is needed because the CompletionResult does not accept an empty string as argument
        if (-Not $Description) {
            $Description = " "
        }
        @{ Name = "$Name"; Description = "$Description" }
    }


    $Space = " "
    if (($Directive -band $ShellCompDirectiveNoSpace) -ne 0 ) {
        # remove the space here
        __${e}_debug "ShellCompDirectiveNoSpace is called"
        $Space = ""
    }

    if ((($Directive -band $ShellCompDirectiveFilterFileExt) -ne 0 ) -or
       (($Directive -band $ShellCompDirectiveFilterDirs) -ne 0 ))  {
        __${e}_debug "ShellCompDirectiveFilterFileExt ShellCompDirectiveFilterDirs are not supported"

        # return here to prevent the completion of the extensions
        return
    }

    $Values = $Values | Where-Object {
        # filter the result
        $_.Name -like "$WordToComplete*"

        # Join the flag back if we have an equal sign flag
        if ( $IsEqualFlag ) {
            __${e}_debug "Join the equal sign flag back to the completion value"
            $_.Name = $Flag + "=" + $_.Name
        }
    }

    # we sort the values in ascending order by name if keep order isn't passed
    if (($Directive -band $ShellCompDirectiveKeepOrder) -eq 0 ) {
        $Values = $Values | Sort-Object -Property Name
    }

    if (($Directive -band $ShellCompDirectiveNoFileComp) -ne 0 ) {
        __${e}_debug "ShellCompDirectiveNoFileComp is called"

        if ($Values.Length -eq 0) {
            # Just print an empty string here so the
            # shell does not start to complete paths.
            # We cannot use CompletionResult here because
            # it does not accept an empty string as argument.
            ""
            return
        }
    }

    # Get the current mode
    $Mode = (Get-PSReadLineKeyHandler | Where-Object { $_.Key -eq "Tab" }).Function
    __${e}_debug "Mode: $Mode"

    $Values | ForEach-Object {

        # store temporary because switch will overwrite $_
        $comp = $_

        # PowerShell supports three different completion modes
        # - TabCompleteNext (default windows style - on each key press the next option is displayed)
        # - Complete (works like bash)
        # - MenuComplete (works like zsh)
        # You set the mode with Set-PSReadLineKeyHandler -Key Tab -Function <mode>

        # CompletionResult Arguments:
        # 1) CompletionText text to be used as the auto completion result
        # 2) ListItemText   text to be displayed in the suggestion list
        # 3) ResultType     type of completion result
        # 4) ToolTip        text for the tooltip with details about the object

        switch ($Mode) {

            # bash like
            "Complete" {

                if ($Values.Length -eq 1) {
                    __${e}_debug "Only one completion left"

                    # insert space after value
                    [System.Management.Automation.CompletionResult]::new($($comp.Name | __${e}_escapeStringWithSpecialChars) + $Space, "$($comp.Name)", 'ParameterValue', "$($comp.Description)")

                } else {
                    # Add the proper number of spaces to align the descriptions
                    while($comp.Name.Length -lt $Longest) {
                        $comp.Name = $comp.Name + " "
                    }

                    # Check for empty description and only add parentheses if needed
                    if ($($comp.Description) -eq " " ) {
                        $Description = ""
                    } else {
                        $Description = "  ($($comp.Description))"
                    }

                    [System.Management.Automation.CompletionResult]::new("$($comp.Name)$Description", "$($comp.Name)$Description", 'ParameterValue', "$($comp.Description)")
                }
             }

            # zsh like
            "MenuComplete" {
                # insert space after value
                # MenuComplete will automatically show the ToolTip of
                # the highlighted value at the bottom of the suggestions.
                [System.Management.Automation.CompletionResult]::new($($comp.Name | __${e}_escapeStringWithSpecialChars) + $Space, "$($comp.Name)", 'ParameterValue', "$($comp.Description)")
            }

            # TabCompleteNext and in case we get something unknown
            Default {
                # Like MenuComplete but we don't want to add a space here because
                # the user need to press space anyway to get the completion.
                # Description will not be shown because that's not possible with TabCompleteNext
                [System.Management.Automation.CompletionResult]::new($($comp.Name | __${e}_escapeStringWithSpecialChars), "$($comp.Name)", 'ParameterValue', "$($comp.Description)")
            }
        }

    }
}

Register-ArgumentCompleter -CommandName '${e}' -ScriptBlock $__${n}CompleterBlock
`}const a={ShellCompDirectiveError:1,ShellCompDirectiveNoSpace:2,ShellCompDirectiveNoFileComp:4,ShellCompDirectiveFilterFileExt:8,ShellCompDirectiveFilterDirs:16,ShellCompDirectiveKeepOrder:32,shellCompDirectiveMaxValue:64,ShellCompDirectiveDefault:0};var o=class{name;variadic;command;handler;constructor(e,t,n,r=!1){this.command=e,this.name=t,this.handler=n,this.variadic=r}},s=class{value;description;command;handler;alias;isBoolean;constructor(e,t,n,r,i,a){this.command=e,this.value=t,this.description=n,this.handler=r,this.alias=i,this.isBoolean=a}},c=class{value;description;options=new Map;arguments=new Map;parent;constructor(e,t){this.value=e,this.description=t}option(e,t,n,r){let i,a,o;typeof n==`function`?(i=n,a=r,o=!1):typeof n==`string`?(i=void 0,a=n,o=!0):(i=void 0,a=void 0,o=!0);let c=new s(this,e,t,i,a,o);return this.options.set(e,c),this}argument(e,t,n=!1){let r=new o(this,e,t,n);return this.arguments.set(e,r),this}},l=class extends c{commands=new Map;completions=[];directive=a.ShellCompDirectiveDefault;constructor(){super(``,``)}command(e,t){let n=new c(e,t);return this.commands.set(e,n),n}stripOptions(e){let t=[],n=0;for(;n<e.length;){let r=e[n];if(r.startsWith(`-`)){n++;let t=!1,i=this.findOption(this,r);if(i)t=i.isBoolean??!1;else for(let[,e]of this.commands){let n=this.findOption(e,r);if(n){t=n.isBoolean??!1;break}}!t&&n<e.length&&!e[n].startsWith(`-`)&&n++}else t.push(r),n++}return t}matchCommand(e){e=this.stripOptions(e);let t=[],n=[],r=null;for(let i=0;i<e.length;i++){let a=e[i];t.push(a);let o=this.commands.get(t.join(` `));if(o)r=o;else{n=e.slice(i,e.length);break}}return[r||this,n]}shouldCompleteFlags(e,t){if(t.startsWith(`-`))return!0;if(e?.startsWith(`-`)){let t=this.findOption(this,e);if(!t){for(let[,n]of this.commands)if(t=this.findOption(n,e),t)break}return!(t&&t.isBoolean)}return!1}shouldCompleteCommands(e){return!e.startsWith(`-`)}handleFlagCompletion(e,t,n,r){let i;if(n.includes(`=`)){let[e]=n.split(`=`);i=e}else if(r?.startsWith(`-`)){let t=this.findOption(e,r);t&&!t.isBoolean&&(i=r)}if(i){let t=this.findOption(e,i);if(t?.handler){let n=[];t.handler.call(t,(e,t)=>n.push({value:e,description:t}),e.options),this.completions=n}return}if(n.startsWith(`-`)){let t=n.startsWith(`-`)&&!n.startsWith(`--`),r=n.replace(/^-+/,``);for(let[i,a]of e.options)t&&a.alias&&`-${a.alias}`.startsWith(n)?this.completions.push({value:`-${a.alias}`,description:a.description}):!t&&i.startsWith(r)&&this.completions.push({value:`--${i}`,description:a.description})}}findOption(e,t){let n=e.options.get(t);if(n||(n=e.options.get(t.replace(/^-+/,``)),n))return n;for(let[n,r]of e.options)if(r.alias&&`-${r.alias}`===t)return r}handleCommandCompletion(e,t){let n=this.stripOptions(e);for(let[e,r]of this.commands){if(e===``)continue;let i=e.split(` `);i.slice(0,n.length).every((e,t)=>e===n[t])&&i[n.length]?.startsWith(t)&&this.completions.push({value:i[n.length],description:r.description})}}handlePositionalCompletion(e,t){let n=e.value.split(` `).length,r=Math.max(0,t.length-n),i=Array.from(e.arguments.entries());if(i.length>0){let t;if(r<i.length){let[e,n]=i[r];t=n}else{let e=i[i.length-1][1];e.variadic&&(t=e)}if(t&&t.handler&&typeof t.handler==`function`){let n=[];t.handler.call(t,(e,t)=>n.push({value:e,description:t}),e.options),this.completions.push(...n)}}}complete(e){this.directive=a.ShellCompDirectiveNoFileComp;let t=new Set;this.completions.filter(e=>t.has(e.value)?!1:(t.add(e.value),!0)).filter(t=>{if(e.includes(`=`)){let[,n]=e.split(`=`);return t.value.startsWith(n)}return t.value.startsWith(e)}).forEach(e=>console.log(`${e.value}\t${e.description??``}`)),console.log(`:${this.directive}`)}parse(e){this.completions=[];let t=e[e.length-1]===``;t&&e.pop();let n=e[e.length-1]||``,r=e.slice(0,-1);t&&(n!==``&&r.push(n),n=``);let[i]=this.matchCommand(r),a=r[r.length-1];if(this.shouldCompleteFlags(a,n))this.handleFlagCompletion(i,r,n,a);else{if(a?.startsWith(`-`)&&n===``&&t){let e=this.findOption(this,a);if(!e){for(let[,t]of this.commands)if(e=this.findOption(t,a),e)break}if(e&&e.isBoolean){this.complete(n);return}}this.shouldCompleteCommands(n)&&this.handleCommandCompletion(r,n),i&&i.arguments.size>0&&this.handlePositionalCompletion(i,r)}this.complete(n)}setup(a,o,s){switch(e(s===`zsh`||s===`bash`||s===`fish`||s===`powershell`,`Unsupported shell`),s){case`zsh`:{let e=t(a,o);console.log(e);break}case`bash`:{let e=n(a,o);console.log(e);break}case`fish`:{let e=r(a,o);console.log(e);break}case`powershell`:{let e=i(a,o);console.log(e);break}}}};const u=new l;function d(e,t,n){u.setup(t,n,e)}var f=u;export{a,i as c,t as d,l as i,r as l,c as n,d as o,s as r,f as s,o as t,n as u};