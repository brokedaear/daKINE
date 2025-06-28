import validator from 'validator'

export function validEmail(email: string): boolean {
  const e = email.trim()
  if (e === '') {
    return false
  }

  console.log(email)

  return validator.isEmail(e)
}

export function validPassword(
  password: string,
  passwordConfirmation: string
): boolean {
  if (password !== passwordConfirmation) {
    return false
  }

  return true
}
