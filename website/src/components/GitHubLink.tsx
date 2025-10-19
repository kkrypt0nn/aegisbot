export enum GitHubType {
  DISCUSSION = "discussions",
  ISSUE = "issues",
  PULL_REQUEST = "pull",
  COMMIT = "commit",
}

type GitHubLinkProps = {
  id: number;
  type: GitHubType;
  hash?: string;
};

export default function GitHubLink({ id, type, hash }: GitHubLinkProps) {
  if (type === GitHubType.COMMIT) {
    return (
      <a href={`https://github.com/kkrypt0nn/aegisbot/${type}/${hash}`}>
        <code>{hash.slice(0, 7)}</code>
      </a>
    );
  }
  return (
    <a href={`https://github.com/kkrypt0nn/aegisbot/${type}/${id}`}>GH-{id}</a>
  );
}
