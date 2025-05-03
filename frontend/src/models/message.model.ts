import { Role } from "@/models/role.model";
export interface Message {
  role: Role;
  content: string | null;
}
