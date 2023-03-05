import Head from 'next/head'
import styles from '@/styles/Home.module.css'
import RedisData from '@/components/redis_data'

export default function Home() {
  return (
    <>
      <Head>
        <title>Price feed</title>
        <meta name="description" content="Price feed monitoring dashboard" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main className={styles.main}>
        {RedisData()}
      </main>
    </>
  )
}
