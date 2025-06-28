import { GitHub } from 'arctic'
import { GithubClientID, GithubClientSecret } from './environment'

export const github = new GitHub(GithubClientID, GithubClientSecret, null)
