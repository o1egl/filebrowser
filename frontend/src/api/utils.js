import { baseURL } from '@/utils/constants'

function getCookies() {
  return document.cookie.split("; ").reduce((c, x) => {
    const splitted = x.split("=");
    c[splitted[0]] = splitted[1];
    return c;
  }, {});
}

export async function fetchURL (url, opts) {
  opts = opts || {}
  opts.headers = opts.headers || {}

  let { headers, ...rest } = opts

  const cookies = getCookies()
  const token = cookies["XSRF-TOKEN"];
  const res = await fetch(pathJoin([baseURL, url]), {
    headers: {
      'X-XSRF-TOKEN': token,
      ...headers
    },
    ...rest
  })

  return res
}

export async function fetchJSON (url, opts) {
  const res = await fetchURL(url, opts)

  if (res.status === 200) {
    return res.json()
  } else {
    throw new Error(res.status)
  }
}

export function removePrefix (url) {
  if (url.startsWith('/files')) {
    url = url.slice(6)
  }

  if (url === '') url = '/'
  if (url[0] !== '/') url = '/' + url
  return url
}

function pathJoin(parts, sep){
  let separator = sep || '/';
  let replace   = new RegExp(separator+'{1,}', 'g');
  return parts.join(separator).replace(replace, separator);
}
