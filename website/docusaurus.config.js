// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const { themes } = require("prism-react-renderer");
const lightCodeTheme = themes.oneLight;
const darkCodeTheme = themes.oneDark;

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: "Aegisbot",
  tagline: "🛡️ Pattern-match your Discord and Twitch defense ",
  // favicon: "assets/purple.png",

  url: "https://aegisbot.krypton.ninja",
  baseUrl: "/",

  organizationName: "kkrypt0nn",
  projectName: "aegisbot",
  trailingSlash: true,

  onBrokenLinks: "throw",
  markdown: {
    hooks: {
      onBrokenMarkdownLinks: "warn",
    },
  },

  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  presets: [
    [
      "classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve("./sidebars.js"),
          editUrl: "https://github.com/kkrypt0nn/aegisbot/tree/main/website",
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      }),
    ],
    [
      "@docusaurus/plugin-sitemap",
      {
        sitemap: {
          lastmod: "date",
          changefreq: "weekly",
          priority: 0.5,
          ignorePatterns: ["/tags/**"],
          filename: "sitemap.xml",
          createSitemapItems: async (params) => {
            const { defaultCreateSitemapItems, ...rest } = params;
            const items = await defaultCreateSitemapItems(rest);
            return items.filter((item) => !item.url.includes("/page/"));
          },
        },
      },
    ],
  ],

  plugins: [
    [
      "@docusaurus/plugin-content-docs",
      /** @type {import('@docusaurus/plugin-content-docs').Options} */
      ({
        id: "community",
        path: "community",
        routeBasePath: "community",
        sidebarPath: require.resolve("./sidebarsCommunity.js"),
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      metadata: [
        {
          name: "keywords",
          content: "aegisbot, discord bot, twitch bot, moderation",
        },
        {
          name: "theme-color",
          content: "#000011",
        },
      ],
      // image: "assets/banner.png",
      colorMode: {
        defaultMode: "dark",
        disableSwitch: false,
        respectPrefersColorScheme: true,
      },
      announcementBar: {
        content:
          "⚠️ Aegisbot is currently under <b>active development and is not yet ready for production use</b>. I strongly advise <b>not</b> using it until a stable release is published. This website is a work in progress as well!",
        backgroundColor: "var(--ifm-navbar-background-color)",
        textColor: "var(--ifm-font-color-base)",
        isCloseable: true,
      },
      navbar: {
        title: "Aegisbot",
        // logo: {
        //   alt: "Aegisbot Logo",
        //   src: "/assets/logo.png",
        // },
        items: [
          {
            type: "doc",
            docId: "README",
            label: "Documentation",
            position: "right",
          },
          {
            to: "/community/code-of-conduct",
            label: "Community",
            position: "right",
          },
          {
            href: "https://github.com/kkrypt0nn/aegisbot",
            "aria-label": "GitHub",
            className: "header-github-link",
            position: "right",
          },
        ],
      },
      footer: {
        style: "dark",
        links: [
          {
            title: "Aegisbot",
            items: [
              {
                label: "Documentation",
                to: "/docs",
              },
              {
                label: "Community",
                to: "/community/code-of-conduct",
              },
            ],
          },
          {
            title: "Community",
            items: [
              {
                label: "Code of Conduct",
                to: "/community/code-of-conduct",
              },
              {
                label: "Contributing Guidelines",
                to: "/community/contributing-guidelines",
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} Aegisbot (Made by Krypton)`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
        additionalLanguages: ["bash"],
      },
    }),
};

module.exports = config;
