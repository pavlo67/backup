package notebook_vews_old

const maxBriefLength = 1000
const GenusKey = "note"
const noteFormID = "note"

//func Genus(ctrlOp controller.Operator, endpoints map[string]config.Endpoint, listeners map[string]config.Listener, pxPreview int) *genera.Genus {
//	return &genera.Genus{
//		ID:         GenusKey,
//		Key:        GenusKey,
//		Name:       "нотатка",
//		NamePlural: "нотатки",
//		Translator: &noteTranslator{
//			ctrlOp: ctrlOp,
//			//host:           host,
//			endpoints:      endpoints,
//			listeners:      listeners,
//			pxPreview:      pxPreview,
//			maxBriefLength: maxBriefLength,
//		},
//	}
//}
//
//const txtNoNote = "Нема такої нотатки"
//
//
//func (gt *noteTranslator) DataFromObject(o *interfaces.Object, context *genera.Context) (map[string]string, interface{}, error) {
//	if o == nil {
//		return nil, nil, nil
//	}
//
//	data := map[string]string{
//		"id":      o.ID,
//		"title":   o.Name,
//		"genus":   o.Genus,
//		"content": o.Content,
//	}
//
//	linksList, err := json.Marshal(o.Links)
//	if err != nil {
//		return nil, nil, errors.Wrapf(err, "can't marshal object.tags: %#v for object.id: %s", o.Links, o.ID)
//	}
//	data["links"] = string(linksList)
//
//	tags := ""
//	filesList := []interfaces.Link{}
//	for _, l := range o.Links {
//		switch l.Type {
//
//		case links.TypeTag:
//			tags += l.Name + "; "
//
//		case files.LinkType:
//			filesList = append(filesList, l)
//		}
//	}
//	data["tags"] = tags
//	if len(filesList) > 0 {
//		files, err := json.Marshal(filesList)
//		if err != nil {
//			log.Println(err)
//		}
//		data["files"] = string(files)
//	}
//
//	if o.UpdatedAt != nil {
//		data["updated_at"] = o.UpdatedAt.Format("02.01.2006 15:04:05")
//	}
//
//	return data, &o, nil
//}
//
//func (gt *noteTranslator) ObjectFromData(userIdentity *confidenter.Identity, oOld *interfaces.Object, data map[string]string, linksList []interfaces.Link) (o *interfaces.Object, index interface{}, context *genera.Context, errs basis.Errors) {
//	if userIdentity == nil {
//		return nil, nil, nil, basis.Errors{confidenter.ErrBadIdentity}
//	}
//
//	var visibility string
//	var rView confidenter.IdentityString
//	var managers = controller.Managers{}
//
//	visibility_ := data["visibility"]
//
//	if visibility_ == string(controller.Anyone) {
//		rView = controller.Anyone
//		visibility = interfaces.Public
//	} else if visibility_ == string(controller.Owner) {
//		rView = userIdentity.String()
//		visibility = interfaces.Private
//	} else if items.ReCommonEdit.MatchString(visibility_) {
//		rView = confidenter.IdentityString(items.ReCommonEdit.ReplaceAllString(visibility_, ""))
//		managers[rights.Change] = rView
//		visibility = interfaces.InGroup
//	} else {
//		rView = confidenter.IdentityString(visibility_)
//		managers[rights.Change] = userIdentity.String()
//		visibility = interfaces.InGroup
//	}
//
//	o = &interfaces.Object{
//		ID:         data["id"],
//		Genus:      data["genus"],
//		Author:     data["author"],
//		Name:       data["title"],
//		Content:    data["content"],
//		Links:      linksList,
//		Visibility: visibility,
//		RView:      rView,
//		Managers:   managers,
//	}
//
//	runes := []rune(o.Content)
//	if len(runes) > gt.maxBriefLength {
//		o.Brief = string(runes[0:gt.maxBriefLength]) + "..."
//	} else {
//		o.Brief = o.Content
//	}
//
//	arr := strings.Split(o.Name, ":")
//	if len(arr) > 1 {
//		o.Links = append(o.Links, interfaces.Link{Type: "author", Name: arr[0]})
//	}
//
//	objects.AddTags(userIdentity, o, data["tags"])
//
//	return o, nil, nil, nil
//}
//
//func (gt *noteTranslator) View(userIdentity *confidenter.Identity, o *interfaces.Object, linkedObjects []interfaces.Object, tab string) map[string]string {
//	if o == nil {
//		return map[string]string{
//			"caput":   "Перегляд",
//			"titulus": "Перегляд",
//			"corpus":  txtNoNote,
//		}
//	}
//
//	var i, htmlIndex, htmlContent, htmlShareTags, htmlShare, linksTitle, htmlLinked string
//	canChange := controller.OneOf(userIdentity, gt.ctrlOp, o.ROwner, o.Managers[rights.Change])
//	canDelete := controller.OneOf(userIdentity, gt.ctrlOp, o.ROwner, o.Managers[rights.Delete])
//
//	if o.Content != "" || canChange || canDelete {
//		htmlIndex = "<tr><td>" + items.HTMLAuthor(userIdentity, o) + "</td></tr>\n"
//	}
//
//	if o.Content != "" {
//		mtext := mt.Read(o.Content)
//		i, htmlContent = mtext.HTML(0, 0)
//		if i != "" {
//			htmlContent = i + "<p>&nbsp;<p>" + htmlContent
//		}
//
//		htmlShareTags = items.HTMLFBTags(o)
//		htmlShare = "<br>&nbsp;<table cellpadding=0 cellspacing=0><tr><td valign=top>" + items.HTMLFBShare(program.Domain()+gt.endpoints["view"].Path(o.ID), o.RView == controller.Anyone) +
//			"<td>&nbsp;<td>" + items.HTMLTwitterShare(o.Name+" "+o.Brief, program.Domain()+gt.endpoints["view"].Path(o.ID)) +
//			"</tr></table>" +
//			items.HTMLFBSDK
//
//		linksTitle = "<p><b>Повʼязані записи, підрозділи</b>\n<p>"
//	}
//
//	if canChange {
//		htmlIndex += `<tr><td>- <a href="` + gt.endpoints["edit"].Path(o.ID) + `">редаґування</a></td></tr>` + "\n"
//	}
//	if canDelete {
//		htmlIndex += `<tr><td>- <a href="#" id="` + gt.listeners["deleteItem"].ID + `">вилучити запис</a>` +
//			`<input type="hidden" id="id" value="` + o.ID + `"></td></tr>` + "\n"
//	}
//
//	htmlIndex += `<tr><td>&nbsp;<br>` + items.HTMLTags(o.Links, o.RView, "", "<br>- ") + "</td></tr>\n"
//
//	data := map[string]string{
//		"caput":      o.Name,
//		"titulus":    o.Name,
//		"share_tags": htmlShareTags,
//		"index":      htmlIndex,
//	}
//
//	if len(linkedObjects) > 0 {
//		htmlLinked = linksTitle + items.HTMLIndex(userIdentity.String(), linkedObjects) + "\n<p>"
//	}
//
//	data["corpus"] = "\n" +
//		htmlContent + "\n<p>" +
//		items.HTMLFiles(o.Links, gt.pxPreview) +
//		htmlShare +
//		htmlLinked
//
//	return data
//}
//
//
//func (gt *noteTranslator) NewItem(user *confidenter.User, o *interfaces.Object, context *genera.Context) map[string]string {
//	data := map[string]string{
//		"genus": GenusKey,
//	}
//
//	rView := user.GetIdentity().String()
//	if o != nil {
//		rView = o.RView
//	}
//
//	return map[string]string{
//		"caput":   "Нова нотатка",
//		"titulus": "Нова нотатка",
//		"corpus":  items.HTMLEdit(user, createFields, noteFormID, data, nil, frontOps, rView, false),
//	}
//
//}
//
//func (gt *noteTranslator) Edit(user *confidenter.User, o *interfaces.Object, context *genera.Context) map[string]string {
//	if o == nil {
//		return map[string]string{
//			"caput":   "Редаґування",
//			"titulus": "Редаґування",
//			"corpus":  txtNoNote,
//		}
//	}
//
//	userIdentity := user.GetIdentity()
//
//	var htmlIndex string
//	if controller.OneOf(userIdentity, gt.ctrlOp, o.ROwner, o.Managers[rights.Delete]) {
//		htmlIndex += `<tr><td>- <a href="#" id="` + gt.listeners["deleteItem"].ID + `">вилучити запис</a>` +
//			`<input type="hidden" id="id" value="` + o.ID + `"></td></tr>` + "\n"
//	}
//
//	htmlIndex += `<tr><td>&nbsp;<br>` + items.HTMLTags(o.Links, o.RView, "", "<br>- ") + "</td></tr>\n"
//
//	responseData := map[string]string{
//		"caput":   o.Name,
//		"titulus": o.Name,
//		"index":   htmlIndex,
//	}
//
//	if !controller.OneOf(userIdentity, gt.ctrlOp, o.ROwner, o.Managers[rights.Change]) {
//		responseData["corpus"] = rights.ErrNoRights.Error()
//		return responseData
//	}
//
//	data, _, err := gt.DataFromObject(o, nil)
//	if err != nil {
//		log.Println(err)
//		responseData["corpus"] = interfaces.ErrCantPerform.Error()
//		return responseData
//
//	}
//
//	publicChanges := o.Managers != nil && o.Managers[rights.Change] == o.RView
//	responseData["corpus"] = items.HTMLEdit(user, updateFields, noteFormID, data, nil, frontOps, o.RView, publicChanges)
//	return responseData
//
//}
