# Disable file completions
complete -c yk -f
# Add completion for setting password
complete -c yk -a set-password -d "Set authentication password"
# Add completion of credentials
complete -c yk -a "(yk list)" -d "Credential to get totp for"

function yk-copy --description "Select a credential using fzf and copy it to clipboard"
    yk (yk list | fzf) | pbcopy
end
