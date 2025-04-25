complete --command yk --no-files
complete --command yk --arguments "(yk list)" --description "Credential to get TOTP for"
complete --command yk --old-option "debug" --description "enable debug logging"
complete --command yk --old-option "version" --description "print version and exit"
complete --command yk --old-option "set-password" --description "prompt for key password and store in system keychain"

function yk-copy --description "Select a credential using fzf and copy it to clipboard"
    set result (yk (yk list | fzf))

    if command -q pbcopy
        echo "$result" | pbcopy
    else if command -q xsel
        echo "$result" | xsel --clipboard
    else if command -q xclip
        echo "$result" | xclip
    else if command -q wl-copy
        echo "$result" | wl-copy
    else
        echo "No clipboard command found. Code is $result"
    end
end
