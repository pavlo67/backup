package files

import "github.com/pavlo67/common/common/joiner"

//  --------------------------------------------------------

const InterfaceKey = joiner.InterfaceKey("files")
const InterfaceCleanerKey = joiner.InterfaceKey("files_cleaner")

const EPSave = "files_save"
const EPRead = "files_read"
const EPRemove = "files_remove"
const EPList = "files_list"
const EPStat = "files_stat"

const SaveHandlerKey = joiner.InterfaceKey(EPSave)
const ReadHandlerKey = joiner.InterfaceKey(EPRead)
const RemoveHandlerKey = joiner.InterfaceKey(EPRemove)
const ListHandlerKey = joiner.InterfaceKey(EPList)
const StatHandlerKey = joiner.InterfaceKey(EPStat)
