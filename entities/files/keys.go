package files

import "github.com/pavlo67/common/common/joiner"

const InterfaceKey = joiner.InterfaceKey("files")
const InterfaceKeyCleaner = joiner.InterfaceKey("files_cleaner")

const EPSave = "files_save"
const EPRead = "files_read"
const EPRemove = "files_remove"
const EPList = "files_list"
const EPStat = "files_stat"

const HandlerSave = joiner.InterfaceKey(EPSave)
const HandlerRead = joiner.InterfaceKey(EPRead)
const HandlerRemove = joiner.InterfaceKey(EPRemove)
const HandlerList = joiner.InterfaceKey(EPList)
const HandlerStat = joiner.InterfaceKey(EPStat)
