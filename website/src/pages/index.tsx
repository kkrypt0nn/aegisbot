import {
  faApple,
  faDocker,
  faLinux,
  faWindows,
} from "@fortawesome/free-brands-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import Layout from "@theme/Layout";

import { config } from "@fortawesome/fontawesome-svg-core";
import "@fortawesome/fontawesome-svg-core/styles.css";
config.autoAddCss = false;

export default function Home() {
  // TODO: Write this entire page better
  return (
    <Layout
      title="Aegisbot"
      description="🛡️ Pattern-match your Discord and Twitch defense "
    >
      <main>
        <div className="container padding-top--md padding-bottom--lg">
          <div style={{ textAlign: "center" }} className="markdown">
            <h1>Aegisbot</h1>
            <p>
              <strong>🛡️ Pattern-match your Discord and Twitch defense </strong>
            </p>
            <p>
              Aegisbot is a novel Discord and Twitch bot with{" "}
              <b>advanced pattern-matching auto-moderation</b>, built on
              concepts inspired by{" "}
              <a href="https://virustotal.github.io/yara/">YARA</a> - not just
              yet another "auto-mod" clone.
            </p>
            <p>
              It allows for fine-grained detection of malicious, spammy, or
              unwanted behavior using customizable matching rules written in a
              simple and common syntax such as YAML or JSON, rather than static
              keyword lists or simplistic triggers.
            </p>
            <p>More features are in development.</p>
            <div className="margin-bottom--xl">
              <a href="/docs" className="button button--primary">
                Documentation
              </a>
            </div>
            <div
              style={{
                display: "flex",
                gap: "20px",
                justifyContent: "center",
                alignItems: "center",
              }}
            >
              <FontAwesomeIcon size="3x" icon={faLinux} />
              <FontAwesomeIcon size="3x" icon={faApple} />
              <FontAwesomeIcon size="3x" icon={faWindows} />
              <FontAwesomeIcon size="3x" icon={faDocker} />
            </div>
          </div>
        </div>
      </main>
    </Layout>
  );
}
