package commands

type CommandInvoker struct {
	command Command
}

func (c *CommandInvoker) SetCommand(command Command) {
	c.command = command
}

func (c *CommandInvoker) ExecuteCommand() error {
	err := c.command.Execute()
	return err
}
