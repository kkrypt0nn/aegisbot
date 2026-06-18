// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  sidebar: [
    {
      type: "category",
      label: "Run",
      items: ["run/docker", "run/from-source"],
      collapsed: false,
    },
    {
      type: "category",
      label: "Rules",
      items: [
        "rules/events",
        "rules/actions",
        "rules/expression",
        "rules/templating",
        "rules/examples",
      ],
      collapsed: false,
    },
    {
      type: "category",
      label: "Objects",
      items: [
        "objects/attachment",
        "objects/channel",
        "objects/member",
        "objects/message",
      ],
      collapsed: true,
    },
  ],
};

module.exports = sidebars;
