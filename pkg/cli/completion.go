package cli

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type completionCommand struct {
	logE   *logrus.Entry
	stdout io.Writer
}

func (cc *completionCommand) command() *cli.Command {
	// https://github.com/lintnet/lintnet/issues/507
	// https://cli.urfave.org/v2/#bash-completion
	return &cli.Command{
		Name:  "completion",
		Usage: "Output shell completion script for bash, zsh, or fish",
		Description: `Output shell completion script for bash, zsh, or fish.
Source the output to enable completion.

e.g.

.bash_profile

if command -v lintnet &> /dev/null; then
	source <(lintnet completion bash)
fi

.zprofile

if command -v lintnet &> /dev/null; then
	source <(lintnet completion zsh)
fi

fish

lintnet completion fish > ~/.config/fish/completions/lintnet.fish
`,
		Subcommands: []*cli.Command{
			{
				Name:   "bash",
				Usage:  "Output shell completion script for bash",
				Action: cc.bashCompletionAction,
			},
			{
				Name:   "zsh",
				Usage:  "Output shell completion script for zsh",
				Action: cc.zshCompletionAction,
			},
			{
				Name:   "fish",
				Usage:  "Output shell completion script for fish",
				Action: cc.fishCompletionAction,
			},
		},
	}
}

func (cc *completionCommand) bashCompletionAction(*cli.Context) error {
	// https://github.com/urfave/cli/blob/main/autocomplete/bash_autocomplete
	// https://github.com/urfave/cli/blob/c3f51bed6fffdf84227c5b59bd3f2e90683314df/autocomplete/bash_autocomplete#L5-L20
	fmt.Fprintln(cc.stdout, `
_cli_bash_autocomplete() {
  if [[ "${COMP_WORDS[0]}" != "source" ]]; then
    local cur opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    if [[ "$cur" == "-"* ]]; then
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} ${cur} --generate-bash-completion )
    else
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
    fi
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
  fi
}

complete -o bashdefault -o default -o nospace -F _cli_bash_autocomplete lintnet`)
	return nil
}

func (cc *completionCommand) zshCompletionAction(*cli.Context) error {
	// https://github.com/urfave/cli/blob/main/autocomplete/zsh_autocomplete
	// https://github.com/urfave/cli/blob/947f9894eef4725a1c15ed75459907b52dde7616/autocomplete/zsh_autocomplete
	fmt.Fprintln(cc.stdout, `#compdef lintnet

_lintnet() {
  local -a opts
  local cur
  cur=${words[-1]}
  if [[ "$cur" == "-"* ]]; then
    opts=("${(@f)$(${words[@]:0:#words[@]-1} ${cur} --generate-bash-completion)}")
  else
    opts=("${(@f)$(${words[@]:0:#words[@]-1} --generate-bash-completion)}")
  fi

  if [[ "${opts[1]}" != "" ]]; then
    _describe 'values' opts
  else
    _files
  fi
}

if [ "$funcstack[1]" = "_lintnet" ]; then
  _lintnet "$@"
else
  compdef _lintnet lintnet
fi`)
	return nil
}

func (cc *completionCommand) fishCompletionAction(c *cli.Context) error {
	s, err := c.App.ToFishCompletion()
	if err != nil {
		return fmt.Errorf("generate fish completion: %w", err)
	}
	fmt.Fprintln(cc.stdout, s)
	return nil
}
