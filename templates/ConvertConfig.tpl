//1.添加以下内容到convlist_petpk_server.xml
<!--{{.CmdArgs.Name}}-->
{{- $CmdArgs := .CmdArgs }}
<CommNode Name="{{.CmdArgs.Name}}Conf">
    {{- print "\n"}}
    {{- $MetaFile := .XmlFile.FileName -}}
    {{- $count := len .XmlFile.Defines -}}
    {{- range $index,$val := .XmlFile.Defines -}}
        {{- if $val.TagOption.Export -}}
{{- print "    "}}<ResNode Name="S_{{$val.Name}}Conf" BinFile="{{$CmdArgs.XmlModulePath}}/{{$val.Name}}Conf.bin" IncludeSheet="{{$val.Desc}}" Meta="{{$val.Name}}" ExcelFile="{{$CmdArgs.ExcelFile}}" Sort="No" BinStyles="0" EntryMapFile="server_conf_meta/{{$MetaFile}}" />
            {{- if lt (add $index 1) $count -}}
    {{- print "\n"}}
            {{- end}}
        {{- end}}
    {{- end}}
</CommNode>

2.使用ResConvert.exe，将xls配置转化为xml

3.执行Xml2CodeGen，生成配置解析代码（参数按需调整，主要预生成文件列表）
./Xml2CodeGen --name {{.CmdArgs.Name}} --author "{{.CmdArgs.Author}}"
{{- range .XmlFile.Defines -}}
    {{- range .LabelTags -}}
{{- print " "}}--label-tag {{.Value}}
    {{- end}}
{{- end -}}
{{- print " "}}--out-path D:\Projects\petcombat-server-proj\src\configlib\src\activity
{{- range .XmlFile.Defines -}}
    {{- if .TagOption.Export -}}
{{- print " "}}{{$CmdArgs.XmlModulePath}}/{{.Name}}Conf.xml
    {{- end}}
{{- end}}