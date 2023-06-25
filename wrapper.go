package main

const shellWrapper=`
# Declare constant for the executable

# Print error message
__ws_err() {
	if [[ $# -eq 0 ]]; then
		echo -n "ws: " 1>&2
		cat <&0 1>&2
	elif [[ $# -eq 1 ]]; then
		echo "ws: $1" 1>&2
	fi
}

# Execute internal workspace command
__ws_internal() {
	__WS_INTERNAL_EXE="ws_internal"
	${__WS_INTERNAL_EXE} "$@"
	return $?
}

# Record the current working directory before changing to workspace
declare WS_PREVIOUS_LOC=""
__ws_push_wd() {
	WS_PREVIOUS_LOC="$(pwd)"
}

# Return to the previously recorded directory
__ws_pop_wd() {
	if [[ -z "${WS_PREVIOUS_LOC}" ]]; then
		ws_err "No working directory was recorded prior to entering the workspace"
		return 1
	fi
	cd "${WS_PREVIOUS_LOC}"
	WS_PREVIOUS_LOC=""
	return 0
}

# Workspace usage instruction
__ws_usage() {
    __ws_err <<_EOF
Usage: ws [command]

ws is a program for managing temporary workspaces.

Commands:
  list|ls          List all available workspaces
  new              Create a new workspace and change into it
  print_current_workspace|pcw
                   Print the current workspace path
  recent|rec       Change to the most recently used workspace
  next|n           Change to the next workspace in the history
  prev|p           Change to the previous workspace in the history
  remove|rm        Remove the current workspace and return to the previous one
  return|ret       Return to the previous workspace
  usage            Show this usage information
_EOF
}

# Workspace function
ws() {
	local cmd="${1:-recent}"

	case "$cmd" in
		"list"|"ls")
			__ws_internal list_workspaces
			;;
		"new")
			__ws_push_wd
			local ws_path="$(__ws_internal create_new_workspace)"
			echo "${ws_path}"
			[[ ! $? -eq 0 ]] && return $?
			cd "${ws_path}"
			;;
		"print_current_workspace"|"pcw")
			__ws_internal print_current_workspace
			[[ ! $? -eq 0 ]] && return $?
			;;
		"recent"|"rec"|"next"|"n"|"prev"|"p")

			local internal_cmd=""
			if [[ "$cmd" == "rec" || "$cmd" == "recent" ]]; then
				__ws_push_wd
				internal_cmd="recent"
			elif [[ "$cmd" == "next" || "$cmd" == "n" ]]; then
				internal_cmd="next"
			elif [[ "$cmd" == "prev" || "$cmd" == "p" ]]; then
				internal_cmd="prev"
			fi

			local ws_path="$(__ws_internal "print_${internal_cmd}_workspace")"
			[[ ! $? -eq 0 ]] && return $?
			cd "${ws_path}"

			;;
		"remove"|"rm")
			local ws_path="$(__ws_internal print_recent_workspace)"
			rm -rf "${ws_path}"
			if [[ -n "${WS_PREVIOUS_LOC}" ]]; then
				__ws_pop_wd || return $?
			else
				local ws_path="$(__ws_internal print_workspace_root)"
				[[ ! $? -eq 0 ]] && return $?
				cd "${ws_path}"
			fi
			;;
		"return"|"ret")
			__ws_pop_wd
			;;
		"usage")
			__ws_usage
			;;
		*)
			__ws_usage
			return 1
			;;
	esac
}
`
