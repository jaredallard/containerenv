package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func generateCompletionBash() string {
	return `
#!/bin/bash

: ${PROG:=$(basename ${BASH_SOURCE})}

_cli_bash_autocomplete() {
		local cur opts base
		COMPREPLY=()
		cur="${COMP_WORDS[COMP_CWORD]}"
		opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
		COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
		return 0
}

complete -F _cli_bash_autocomplete containerenv

unset PROG
	`
}

func generateCompletionZsh() string {
	return `
_cli_zsh_autocomplete() {

	local -a opts
	opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[@]:0:#words[@]-1} --generate-bash-completion)}")

	_describe 'values' opts

	return
}

compdef _cli_zsh_autocomplete containerenv
	`
}

func generateCompletion(c *cli.Context) error {
	shell := c.Args().First()

	if shell == "" {
		return fmt.Errorf("Please provide a shell")
	}

	v := ""
	switch shell {
	case "zsh":
		v = generateCompletionZsh()
	case "bash":
		v = generateCompletionBash()
	default:
		return fmt.Errorf("Shell '%s' not supported", shell)
	}

	fmt.Println(v)

	return nil
}
