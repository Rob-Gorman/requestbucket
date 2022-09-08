function urlGenerator() {
  // define ranges for valid characters
  // numbers 48-57
  // capital letters 65-90
  // lowercase letters 97-122
  // pick out 9 and return the string
  const chars = [0,1,2,3,4,5,6,7,8,9,'a','b','c','d','e','f','g','h','i','j','k','l','m','n','o','p','q','r','s','t','u','v','w','x','y','z','A','B','C','D','E','F','G','H','I','J','K','L','M','N','O','P','Q','R','S','T','U','V','W','X','Y','Z']

  let url = '';

  for (let c = 0; c < 9; c++) {
    let rand = Math.floor(Math.random() * 62);
    url += chars[rand];
  }

  return url;
}

module.exports = { urlGenerator };