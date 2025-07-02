<!--{{.CmdArgs.Name}}-->
{{- $ExcelFile := .CmdArgs.ExcelFile }}
<CommNode Name="{{.CmdArgs.Name}}Conf">
    {{- print "\n"}}
    {{- $MetaFile := .XmlFile.FileName -}}
    {{- range .XmlFile.Defines -}}
        {{- if .TagOption.Export -}}
{{- print "    "}}<ResNode Name="S_{{.Name}}Conf" BinFile="xls/server/{{.Name}}Conf.bin" IncludeSheet="{{.Desc}}" Meta="{{.Name}}" ExcelFile="{{$ExcelFile}}" Sort="No" BinStyles="0" EntryMapFile="server_conf_meta/{{$MetaFile}}" />
    {{- print "\n"}}
        {{- end}}
    {{- end}}
</CommNode>