import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: 'General purpose linter',
    description: (
      <>
        lintnet is a general purpose linter that can be used for any kind of files.
        lintnet is a command line tool like <a href="https://www.conftest.dev/" target="_blank">Conftest</a>.
        lintnet doesn't have any built-in lint rules, but you can define lint rules yourself freely.
        lintnet supports <a href="/docs/supported-data-format/">various file formats such as JSON, YAML, TOML, HCL, and so on</a>.
      </>
    ),
  },
  {
    title: 'Powered by Jsonnet',
    description: (
      <>
        lintnet uses <a href="https://jsonnet.org/" target="_blank">Jsonnet</a> for configuration files and lint rules.
        Jsonnet is a well-known, simple, and powerful configuration language.
        Jsonnet is a superset of JSON, and can generate JSON dynamically.
        Jsonnet is easy to understand and the learning cost is low.
        Jsonnet is widely used, so you can utilize the knowledge and experience of Jsonnet in not only lintnet but also other projects.
      </>
    ),
  },
  {
    title: 'Secure',
    description: (
      <>
        Jsonnet restricts access to filesystem and network so it's secure compared with common programming languages such as Python
      </>
    ),
  },
  {
    title: 'Shareable',
    description: (
      <>
        lintnet supports sharing lint rules as <a href="/docs/module/">Modules</a>.
        You can utilize third party lint rules, reuse your lint rules in multiple projects, and distribute lint rules as OSS and in-house libraries
      </>
    ),
  },
];

function Feature({ title, description }) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
