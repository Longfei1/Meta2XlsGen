//1.添加以下内容到convlist_petpk_server.xml
<!--{{.CmdArgs.Name}}-->
{{- $CmdArgs := .CmdArgs }}
<CommNode Name="{{.CmdArgs.Name}}Conf">
    {{- print "\n"}}
    {{- $MetaFile := .XmlFile.FileName -}}
    {{- range .XmlFile.Defines -}}
        {{- if .TagOption.Export -}}
{{- print "    "}}<ResNode Name="S_{{.Name}}Conf" BinFile="{{$CmdArgs.ExportPath}}/{{.Name}}Conf.bin" IncludeSheet="{{.Desc}}" Meta="{{.Name}}" ExcelFile="{{$CmdArgs.ExcelFile}}" Sort="No" BinStyles="0" EntryMapFile="server_conf_meta/{{$MetaFile}}" />
    {{- print "\n"}}
        {{- end}}
    {{- end}}
</CommNode>

2.使用ResConvert.exe，将xls配置转化为xml

3.执行Xml2CodeGen，生成配置解析代码（参数按需调整，主要预生成文件列表）
./Xml2CodeGen --name {{.CmdArgs.Name}} --author "{{.CmdArgs.Author}}"
{{- range .XmlFile.Defines -}}
    {{- if .TagOption.Export -}}
{{- print " "}}{{$CmdArgs.ExportPath}}/{{.Name}}Conf.xml
    {{- end}}
{{- end}}