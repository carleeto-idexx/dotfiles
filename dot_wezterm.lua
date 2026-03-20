-- Pull in the wezterm API
local wezterm = require("wezterm")

-- This will hold the configuration.
local config = wezterm.config_builder()

-- This is where you actually apply your config choices.

-- For example, changing the initial geometry for new windows:
config.initial_cols = 120
config.initial_rows = 28

-- or, changing the font size and color scheme.
config.font_size = 12
config.color_scheme = "Dracula (Official)"

config.window_decorations = "RESIZE"
config.enable_tab_bar = false
config.window_close_confirmation = "NeverPrompt"

config.keys = {
	{
		key = "d",
		mods = "CMD",
		action = wezterm.action.SplitHorizontal({ domain = "CurrentPaneDomain" }),
	},
	{
		key = "d",
		mods = "CMD|SHIFT",
		action = wezterm.action.SplitVertical({ domain = "CurrentPaneDomain" }),
	},
	{
		key = "k",
		mods = "CMD",
		action = wezterm.action.SendString("clear\n"),
	},
	{
		key = "s",
		mods = "CTRL|SHIFT",
		action = wezterm.action_callback(function(window, pane)
			local tabs = {}
			for _, tab_info in ipairs(window:mux_window():tabs_with_info()) do
				local tab_pane = tab_info.tab:active_pane()
				if tab_pane:get_domain_name() == "unix" then
					local cwd = tab_pane:get_current_working_dir()
					local dir = cwd and cwd.file_path or ""
					local process = tab_pane:get_foreground_process_name() or ""
					process = process:match("([^/]+)$") or ""
					if process == "" then
						process = tab_pane:get_title() or ""
					end
					local label = string.format(
						"%d: %s (%s)",
						tab_info.index + 1,
						process,
						dir
					)
					table.insert(tabs, {
						id = tostring(tab_info.tab:tab_id()),
						label = label,
					})
				end
			end
			window:perform_action(
				wezterm.action.InputSelector({
					title = "Unix Domain Tabs",
					choices = tabs,
					action = wezterm.action_callback(function(win, _, id)
						if id then
							for _, tab_info in ipairs(win:mux_window():tabs_with_info()) do
								if tostring(tab_info.tab:tab_id()) == id then
									tab_info.tab:activate()
									break
								end
							end
						end
					end),
				}),
				pane
			)
		end),
	},
}

-- Finally, return the configuration to wezterm:
return config
