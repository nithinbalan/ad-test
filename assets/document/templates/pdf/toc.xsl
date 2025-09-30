<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="2.0"
                xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                xmlns:outline="http://wkhtmltopdf.org/outline"
                xmlns="http://www.w3.org/1999/xhtml">
  <xsl:output doctype-public="-//W3C//DTD XHTML 1.0 Strict//EN"
              doctype-system="http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd"
              indent="yes" />
  <xsl:template match="outline:outline">
    <html>
      <head>
        <title>Table of Contents</title>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
        <style>
          h1 {
            text-align: left;
            font-size: 18px;
            font-family: Lexend Deca, sans-serif;
            color: #050505; 
            font-weight: 400;
          }
          div {border-bottom: 1px dashed rgb(200,200,200);}
          span {float: right;}
          li {
            list-style: none;
            line-height: 1.6;
            }
          ul {
            font-size: 14px;  /* same size for all */
            font-family: Lexend Deca, sans-serif;
            padding-left: 0em;
            line-height: 1.6;
          }
          /* remove ul ul or set it to same size */
          ul ul {font-size: 14px; padding-left: 1em;}
          ul {padding-left: 0em;}
          ul ul {padding-left: 0em;}
          a {text-decoration:none; color: black;}
        </style>
      </head>
      <body>
        <h1>Table of Contents</h1>
        <ul><xsl:apply-templates select="outline:item/outline:item[@title != 'Table of Contents']"/></ul>
      </body>
    </html>
  </xsl:template>
  <xsl:template match="outline:item">
    <li>
      <xsl:if test="@title!='' and @link != ''">
        <div>
          <a>
            <xsl:if test="@link">
              <xsl:attribute name="href"><xsl:value-of select="@link"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@backLink">
              <xsl:attribute name="name"><xsl:value-of select="@backLink"/></xsl:attribute>
            </xsl:if>
            <xsl:value-of select="@title" /> 
          </a>
          <span> <xsl:value-of select="@page" /> </span>
        </div>
      </xsl:if>
      <ul>
        <xsl:apply-templates select="outline:item[contains(@class,'toc-item')]"/>
      </ul>
    </li>
  </xsl:template>
</xsl:stylesheet>
