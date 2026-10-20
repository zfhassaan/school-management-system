export namespace academics {
	
	export class AcademicTerm {
	    ID: string;
	    AcademicYearID: string;
	    Sequence: number;
	    Name: string;
	    StartsOn: shared.Date;
	    EndsOn: shared.Date;
	
	    static createFrom(source: any = {}) {
	        return new AcademicTerm(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.AcademicYearID = source["AcademicYearID"];
	        this.Sequence = source["Sequence"];
	        this.Name = source["Name"];
	        this.StartsOn = this.convertValues(source["StartsOn"], shared.Date);
	        this.EndsOn = this.convertValues(source["EndsOn"], shared.Date);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AcademicTermInput {
	    sequence: number;
	    name: string;
	    startsOn: shared.Date;
	    endsOn: shared.Date;
	
	    static createFrom(source: any = {}) {
	        return new AcademicTermInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sequence = source["sequence"];
	        this.name = source["name"];
	        this.startsOn = this.convertValues(source["startsOn"], shared.Date);
	        this.endsOn = this.convertValues(source["endsOn"], shared.Date);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AcademicYear {
	    ID: string;
	    SchoolID: string;
	    Code: string;
	    Name: string;
	    StartsOn: shared.Date;
	    EndsOn: shared.Date;
	    Status: string;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	
	    static createFrom(source: any = {}) {
	        return new AcademicYear(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.SchoolID = source["SchoolID"];
	        this.Code = source["Code"];
	        this.Name = source["Name"];
	        this.StartsOn = this.convertValues(source["StartsOn"], shared.Date);
	        this.EndsOn = this.convertValues(source["EndsOn"], shared.Date);
	        this.Status = source["Status"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Campus {
	    ID: string;
	    SchoolID: string;
	    Code: string;
	    Name: string;
	    Address: string;
	    IsActive: boolean;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	
	    static createFrom(source: any = {}) {
	        return new Campus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.SchoolID = source["SchoolID"];
	        this.Code = source["Code"];
	        this.Name = source["Name"];
	        this.Address = source["Address"];
	        this.IsActive = source["IsActive"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Class {
	    ID: string;
	    SchoolID: string;
	    Code: string;
	    Name: string;
	    LevelOrder: number;
	    NextClassID?: string;
	    IsActive: boolean;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	
	    static createFrom(source: any = {}) {
	        return new Class(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.SchoolID = source["SchoolID"];
	        this.Code = source["Code"];
	        this.Name = source["Name"];
	        this.LevelOrder = source["LevelOrder"];
	        this.NextClassID = source["NextClassID"];
	        this.IsActive = source["IsActive"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CreateAcademicYearInput {
	    name: string;
	    code: string;
	    startsOn: shared.Date;
	    endsOn: shared.Date;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateAcademicYearInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.code = source["code"];
	        this.startsOn = this.convertValues(source["startsOn"], shared.Date);
	        this.endsOn = this.convertValues(source["endsOn"], shared.Date);
	        this.status = source["status"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CreateCampusInput {
	    name: string;
	    code: string;
	    address: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateCampusInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.code = source["code"];
	        this.address = source["address"];
	    }
	}
	export class CreateClassInput {
	    name: string;
	    code: string;
	    levelOrder: number;
	    nextClassId?: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateClassInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.code = source["code"];
	        this.levelOrder = source["levelOrder"];
	        this.nextClassId = source["nextClassId"];
	    }
	}
	export class CreateSectionInput {
	    classId: string;
	    academicYearId: string;
	    campusId?: string;
	    name: string;
	    capacity?: number;
	    classTeacherUserId?: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateSectionInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.classId = source["classId"];
	        this.academicYearId = source["academicYearId"];
	        this.campusId = source["campusId"];
	        this.name = source["name"];
	        this.capacity = source["capacity"];
	        this.classTeacherUserId = source["classTeacherUserId"];
	    }
	}
	export class Section {
	    ID: string;
	    ClassID: string;
	    AcademicYearID: string;
	    CampusID?: string;
	    Name: string;
	    Capacity?: number;
	    ClassTeacherUserID?: string;
	    IsActive: boolean;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	
	    static createFrom(source: any = {}) {
	        return new Section(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.ClassID = source["ClassID"];
	        this.AcademicYearID = source["AcademicYearID"];
	        this.CampusID = source["CampusID"];
	        this.Name = source["Name"];
	        this.Capacity = source["Capacity"];
	        this.ClassTeacherUserID = source["ClassTeacherUserID"];
	        this.IsActive = source["IsActive"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SetTermsInput {
	    academicYearId: string;
	    terms: AcademicTermInput[];
	
	    static createFrom(source: any = {}) {
	        return new SetTermsInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.academicYearId = source["academicYearId"];
	        this.terms = this.convertValues(source["terms"], AcademicTermInput);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateAcademicYearInput {
	    id: string;
	    name: string;
	    startsOn: shared.Date;
	    endsOn: shared.Date;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateAcademicYearInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.startsOn = this.convertValues(source["startsOn"], shared.Date);
	        this.endsOn = this.convertValues(source["endsOn"], shared.Date);
	        this.status = source["status"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateClassInput {
	    id: string;
	    name: string;
	    code: string;
	    levelOrder: number;
	    nextClassId?: string;
	    isActive: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UpdateClassInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.code = source["code"];
	        this.levelOrder = source["levelOrder"];
	        this.nextClassId = source["nextClassId"];
	        this.isActive = source["isActive"];
	    }
	}
	export class UpdateSectionInput {
	    id: string;
	    name: string;
	    capacity?: number;
	    classTeacherUserId?: string;
	    isActive: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UpdateSectionInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.capacity = source["capacity"];
	        this.classTeacherUserId = source["classTeacherUserId"];
	        this.isActive = source["isActive"];
	    }
	}

}

export namespace admissions {
	
	export class Application {
	    ID: string;
	    SchoolID: string;
	    ApplicationNo: string;
	    FormVersionID: string;
	    AcademicYearID: string;
	    AppliedClassID: string;
	    Status: string;
	    ApplicantFirstName: string;
	    ApplicantLastName: string;
	    ApplicantDOB?: shared.Date;
	    ApplicantGender: string;
	    GuardianPhoneNormalized: string;
	    StudentID?: string;
	    Response: Record<string, any>;
	    // Go type: time
	    SubmittedAt?: any;
	    // Go type: time
	    DecidedAt?: any;
	    DecidedBy?: string;
	    DecisionReason: string;
	    CreatedBy: string;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	    ClassName: string;
	    YearCode: string;
	
	    static createFrom(source: any = {}) {
	        return new Application(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.SchoolID = source["SchoolID"];
	        this.ApplicationNo = source["ApplicationNo"];
	        this.FormVersionID = source["FormVersionID"];
	        this.AcademicYearID = source["AcademicYearID"];
	        this.AppliedClassID = source["AppliedClassID"];
	        this.Status = source["Status"];
	        this.ApplicantFirstName = source["ApplicantFirstName"];
	        this.ApplicantLastName = source["ApplicantLastName"];
	        this.ApplicantDOB = this.convertValues(source["ApplicantDOB"], shared.Date);
	        this.ApplicantGender = source["ApplicantGender"];
	        this.GuardianPhoneNormalized = source["GuardianPhoneNormalized"];
	        this.StudentID = source["StudentID"];
	        this.Response = source["Response"];
	        this.SubmittedAt = this.convertValues(source["SubmittedAt"], null);
	        this.DecidedAt = this.convertValues(source["DecidedAt"], null);
	        this.DecidedBy = source["DecidedBy"];
	        this.DecisionReason = source["DecisionReason"];
	        this.CreatedBy = source["CreatedBy"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	        this.ClassName = source["ClassName"];
	        this.YearCode = source["YearCode"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CreateFormDefInput {
	    code: string;
	    name: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateFormDefInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.name = source["name"];
	        this.description = source["description"];
	    }
	}
	export class Visibility {
	    fieldId: string;
	    equals: string;
	
	    static createFrom(source: any = {}) {
	        return new Visibility(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fieldId = source["fieldId"];
	        this.equals = source["equals"];
	    }
	}
	export class FormField {
	    id: string;
	    label: string;
	    type: string;
	    placeholder?: string;
	    required: boolean;
	    options?: string[];
	    dependsOn?: Visibility;
	
	    static createFrom(source: any = {}) {
	        return new FormField(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	        this.type = source["type"];
	        this.placeholder = source["placeholder"];
	        this.required = source["required"];
	        this.options = source["options"];
	        this.dependsOn = this.convertValues(source["dependsOn"], Visibility);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FormSection {
	    id: string;
	    title: string;
	    description?: string;
	    fields: FormField[];
	
	    static createFrom(source: any = {}) {
	        return new FormSection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.fields = this.convertValues(source["fields"], FormField);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FormSchema {
	    title: string;
	    description?: string;
	    sections: FormSection[];
	    requiredDocuments?: string[];
	
	    static createFrom(source: any = {}) {
	        return new FormSchema(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.description = source["description"];
	        this.sections = this.convertValues(source["sections"], FormSection);
	        this.requiredDocuments = source["requiredDocuments"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CreateFormVersionInput {
	    formDefinitionId: string;
	    schema: FormSchema;
	
	    static createFrom(source: any = {}) {
	        return new CreateFormVersionInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.formDefinitionId = source["formDefinitionId"];
	        this.schema = this.convertValues(source["schema"], FormSchema);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EnrollApplicationInput {
	    applicationId: string;
	    sectionId?: string;
	    rollNo?: string;
	    admissionNo: string;
	    admissionDate: shared.Date;
	
	    static createFrom(source: any = {}) {
	        return new EnrollApplicationInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.applicationId = source["applicationId"];
	        this.sectionId = source["sectionId"];
	        this.rollNo = source["rollNo"];
	        this.admissionNo = source["admissionNo"];
	        this.admissionDate = this.convertValues(source["admissionDate"], shared.Date);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FormDefinition {
	    ID: string;
	    SchoolID: string;
	    Code: string;
	    Name: string;
	    Description: string;
	    IsActive: boolean;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	
	    static createFrom(source: any = {}) {
	        return new FormDefinition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.SchoolID = source["SchoolID"];
	        this.Code = source["Code"];
	        this.Name = source["Name"];
	        this.Description = source["Description"];
	        this.IsActive = source["IsActive"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	export class FormVersion {
	    ID: string;
	    FormDefinitionID: string;
	    VersionNo: number;
	    Status: string;
	    Schema: FormSchema;
	    SchemaHash: string;
	    // Go type: time
	    PublishedAt?: any;
	    PublishedBy?: string;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	
	    static createFrom(source: any = {}) {
	        return new FormVersion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.FormDefinitionID = source["FormDefinitionID"];
	        this.VersionNo = source["VersionNo"];
	        this.Status = source["Status"];
	        this.Schema = this.convertValues(source["Schema"], FormSchema);
	        this.SchemaHash = source["SchemaHash"];
	        this.PublishedAt = this.convertValues(source["PublishedAt"], null);
	        this.PublishedBy = source["PublishedBy"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PublishFormVersionInput {
	    versionId: string;
	
	    static createFrom(source: any = {}) {
	        return new PublishFormVersionInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.versionId = source["versionId"];
	    }
	}
	export class SubmitApplicationInput {
	    formVersionId: string;
	    academicYearId: string;
	    appliedClassId: string;
	    applicantFirstName: string;
	    applicantLastName: string;
	    applicantDob?: shared.Date;
	    applicantGender: string;
	    guardianPhoneNormalized: string;
	    guardianName: string;
	    guardianEmail: string;
	    response: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new SubmitApplicationInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.formVersionId = source["formVersionId"];
	        this.academicYearId = source["academicYearId"];
	        this.appliedClassId = source["appliedClassId"];
	        this.applicantFirstName = source["applicantFirstName"];
	        this.applicantLastName = source["applicantLastName"];
	        this.applicantDob = this.convertValues(source["applicantDob"], shared.Date);
	        this.applicantGender = source["applicantGender"];
	        this.guardianPhoneNormalized = source["guardianPhoneNormalized"];
	        this.guardianName = source["guardianName"];
	        this.guardianEmail = source["guardianEmail"];
	        this.response = source["response"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateApplicationStatusInput {
	    applicationId: string;
	    newStatus: string;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateApplicationStatusInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.applicationId = source["applicationId"];
	        this.newStatus = source["newStatus"];
	        this.reason = source["reason"];
	    }
	}

}

export namespace audit {
	
	export class AuditLogDTO {
	    id: string;
	    occurredAt: string;
	    actorUserId?: string;
	    deviceId: string;
	    action: string;
	    entityType: string;
	    entityId: string;
	    correlationId: string;
	    reason: string;
	    before?: any;
	    after?: any;
	    source: string;
	
	    static createFrom(source: any = {}) {
	        return new AuditLogDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.occurredAt = source["occurredAt"];
	        this.actorUserId = source["actorUserId"];
	        this.deviceId = source["deviceId"];
	        this.action = source["action"];
	        this.entityType = source["entityType"];
	        this.entityId = source["entityId"];
	        this.correlationId = source["correlationId"];
	        this.reason = source["reason"];
	        this.before = source["before"];
	        this.after = source["after"];
	        this.source = source["source"];
	    }
	}
	export class Filter {
	    EntityType: string;
	    EntityID: string;
	    ActorID?: string;
	    // Go type: time
	    From?: any;
	    // Go type: time
	    To?: any;
	    Cursor: string;
	    Limit: number;
	
	    static createFrom(source: any = {}) {
	        return new Filter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.EntityType = source["EntityType"];
	        this.EntityID = source["EntityID"];
	        this.ActorID = source["ActorID"];
	        this.From = this.convertValues(source["From"], null);
	        this.To = this.convertValues(source["To"], null);
	        this.Cursor = source["Cursor"];
	        this.Limit = source["Limit"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace auth {
	
	export class ChangePasswordInput {
	    oldPassword: string;
	    newPassword: string;
	
	    static createFrom(source: any = {}) {
	        return new ChangePasswordInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.oldPassword = source["oldPassword"];
	        this.newPassword = source["newPassword"];
	    }
	}
	export class FirstRunSetupInput {
	    schoolName: string;
	    schoolCode: string;
	    timezone: string;
	    currencyCode: string;
	    currencyMinorUnits: number;
	    adminUsername: string;
	    adminDisplayName: string;
	    adminEmail: string;
	    adminPassword: string;
	
	    static createFrom(source: any = {}) {
	        return new FirstRunSetupInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.schoolName = source["schoolName"];
	        this.schoolCode = source["schoolCode"];
	        this.timezone = source["timezone"];
	        this.currencyCode = source["currencyCode"];
	        this.currencyMinorUnits = source["currencyMinorUnits"];
	        this.adminUsername = source["adminUsername"];
	        this.adminDisplayName = source["adminDisplayName"];
	        this.adminEmail = source["adminEmail"];
	        this.adminPassword = source["adminPassword"];
	    }
	}
	export class FirstRunSetupResult {
	    schoolId: string;
	    userId: string;
	    recoveryCodes: string[];
	
	    static createFrom(source: any = {}) {
	        return new FirstRunSetupResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.schoolId = source["schoolId"];
	        this.userId = source["userId"];
	        this.recoveryCodes = source["recoveryCodes"];
	    }
	}
	export class LoginInput {
	    username: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new LoginInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.password = source["password"];
	    }
	}
	export class SchoolDTO {
	    id: string;
	    name: string;
	    code: string;
	    timezone: string;
	    currencyCode: string;
	    currencyMinorUnits: number;
	    dateFormat: string;
	
	    static createFrom(source: any = {}) {
	        return new SchoolDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.code = source["code"];
	        this.timezone = source["timezone"];
	        this.currencyCode = source["currencyCode"];
	        this.currencyMinorUnits = source["currencyMinorUnits"];
	        this.dateFormat = source["dateFormat"];
	    }
	}
	export class UserDTO {
	    id: string;
	    username: string;
	    displayName: string;
	    email: string;
	    roles: string[];
	    permissions: string[];
	    mustChangePassword: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UserDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	        this.displayName = source["displayName"];
	        this.email = source["email"];
	        this.roles = source["roles"];
	        this.permissions = source["permissions"];
	        this.mustChangePassword = source["mustChangePassword"];
	    }
	}
	export class LoginResult {
	    sessionToken: string;
	    user: UserDTO;
	    school: SchoolDTO;
	
	    static createFrom(source: any = {}) {
	        return new LoginResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sessionToken = source["sessionToken"];
	        this.user = this.convertValues(source["user"], UserDTO);
	        this.school = this.convertValues(source["school"], SchoolDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class SetupStatusDTO {
	    isInitialized: boolean;
	    schoolName?: string;
	
	    static createFrom(source: any = {}) {
	        return new SetupStatusDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isInitialized = source["isInitialized"];
	        this.schoolName = source["schoolName"];
	    }
	}

}

export namespace documents {
	
	export class Document {
	    ID: string;
	    StudentID?: string;
	    ApplicationID?: string;
	    DocumentTypeCode: string;
	    FileName: string;
	    MimeType: string;
	    SizeBytes: number;
	    SHA256: string;
	    StorageKey: string;
	    VersionNo: number;
	    UploadedBy: string;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	
	    static createFrom(source: any = {}) {
	        return new Document(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.StudentID = source["StudentID"];
	        this.ApplicationID = source["ApplicationID"];
	        this.DocumentTypeCode = source["DocumentTypeCode"];
	        this.FileName = source["FileName"];
	        this.MimeType = source["MimeType"];
	        this.SizeBytes = source["SizeBytes"];
	        this.SHA256 = source["SHA256"];
	        this.StorageKey = source["StorageKey"];
	        this.VersionNo = source["VersionNo"];
	        this.UploadedBy = source["UploadedBy"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace dto {
	
	export class ErrorDTO {
	    code: string;
	    message: string;
	    fields?: Record<string, string>;
	    correlationId?: string;
	
	    static createFrom(source: any = {}) {
	        return new ErrorDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	        this.fields = source["fields"];
	        this.correlationId = source["correlationId"];
	    }
	}
	export class Response__school_management_internal_application_commands_auth_FirstRunSetupResult_ {
	    success: boolean;
	    data?: auth.FirstRunSetupResult;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_application_commands_auth_FirstRunSetupResult_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], auth.FirstRunSetupResult);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_application_commands_auth_LoginResult_ {
	    success: boolean;
	    data?: auth.LoginResult;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_application_commands_auth_LoginResult_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], auth.LoginResult);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_application_commands_payments_RecordPaymentResult_ {
	    success: boolean;
	    data?: payments.RecordPaymentResult;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_application_commands_payments_RecordPaymentResult_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], payments.RecordPaymentResult);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_application_queries_auth_SetupStatusDTO_ {
	    success: boolean;
	    data?: auth.SetupStatusDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_application_queries_auth_SetupStatusDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], auth.SetupStatusDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_application_queries_fees_FeeStructureVersionDetail_ {
	    success: boolean;
	    data?: fees.FeeStructureVersionDetail;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_application_queries_fees_FeeStructureVersionDetail_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], fees.FeeStructureVersionDetail);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_application_queries_fees_InvoiceDetail_ {
	    success: boolean;
	    data?: fees.InvoiceDetail;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_application_queries_fees_InvoiceDetail_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], fees.InvoiceDetail);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_application_queries_payments_PaymentDetail_ {
	    success: boolean;
	    data?: payments.PaymentDetail;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_application_queries_payments_PaymentDetail_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], payments.PaymentDetail);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_application_queries_payments_StudentCreditBalanceInfo_ {
	    success: boolean;
	    data?: payments.StudentCreditBalanceInfo;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_application_queries_payments_StudentCreditBalanceInfo_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], payments.StudentCreditBalanceInfo);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_application_queries_roles_RoleDTO_ {
	    success: boolean;
	    data?: roles.RoleDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_application_queries_roles_RoleDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], roles.RoleDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_application_queries_users_UserListItemDTO_ {
	    success: boolean;
	    data?: users.UserListItemDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_application_queries_users_UserListItemDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], users.UserListItemDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_domain_academics_AcademicYear_ {
	    success: boolean;
	    data?: academics.AcademicYear;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_domain_academics_AcademicYear_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], academics.AcademicYear);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_domain_admissions_Application_ {
	    success: boolean;
	    data?: admissions.Application;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_domain_admissions_Application_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], admissions.Application);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_domain_admissions_FormVersion_ {
	    success: boolean;
	    data?: admissions.FormVersion;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_domain_admissions_FormVersion_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], admissions.FormVersion);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_domain_guardians_Guardian_ {
	    success: boolean;
	    data?: guardians.Guardian;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_domain_guardians_Guardian_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], guardians.Guardian);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_domain_payments_FinancialAdjustment_ {
	    success: boolean;
	    data?: payments.FinancialAdjustment;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_domain_payments_FinancialAdjustment_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], payments.FinancialAdjustment);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_domain_school_School_ {
	    success: boolean;
	    data?: school.School;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_domain_school_School_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], school.School);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_domain_school_Settings_ {
	    success: boolean;
	    data?: school.Settings;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_domain_school_Settings_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], school.Settings);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_domain_students_Enrollment_ {
	    success: boolean;
	    data?: students.Enrollment;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_domain_students_Enrollment_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], students.Enrollment);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_domain_students_PromotionBatch_ {
	    success: boolean;
	    data?: students.PromotionBatch;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_domain_students_PromotionBatch_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], students.PromotionBatch);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_domain_students_Student_ {
	    success: boolean;
	    data?: students.Student;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_domain_students_Student_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], students.Student);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response__school_management_internal_transport_wails_DocumentContentResponse_ {
	    success: boolean;
	    data?: wails.DocumentContentResponse;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response__school_management_internal_transport_wails_DocumentContentResponse_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], wails.DocumentContentResponse);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response____school_management_internal_domain_payments_FinancialAdjustment_ {
	    success: boolean;
	    data?: payments.FinancialAdjustment[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response____school_management_internal_domain_payments_FinancialAdjustment_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], payments.FinancialAdjustment);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response____school_management_internal_domain_payments_PaymentWithDetails_ {
	    success: boolean;
	    data?: payments.PaymentWithDetails[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response____school_management_internal_domain_payments_PaymentWithDetails_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], payments.PaymentWithDetails);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_application_queries_audit_AuditLogDTO_ {
	    success: boolean;
	    data?: audit.AuditLogDTO[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_application_queries_audit_AuditLogDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], audit.AuditLogDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_application_queries_roles_PermissionInfoDTO_ {
	    success: boolean;
	    data?: roles.PermissionInfoDTO[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_application_queries_roles_PermissionInfoDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], roles.PermissionInfoDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_application_queries_roles_RoleDTO_ {
	    success: boolean;
	    data?: roles.RoleDTO[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_application_queries_roles_RoleDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], roles.RoleDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_application_queries_users_UserListItemDTO_ {
	    success: boolean;
	    data?: users.UserListItemDTO[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_application_queries_users_UserListItemDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], users.UserListItemDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_academics_AcademicTerm_ {
	    success: boolean;
	    data?: academics.AcademicTerm[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_academics_AcademicTerm_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], academics.AcademicTerm);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_academics_AcademicYear_ {
	    success: boolean;
	    data?: academics.AcademicYear[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_academics_AcademicYear_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], academics.AcademicYear);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_academics_Campus_ {
	    success: boolean;
	    data?: academics.Campus[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_academics_Campus_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], academics.Campus);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_academics_Class_ {
	    success: boolean;
	    data?: academics.Class[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_academics_Class_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], academics.Class);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_academics_Section_ {
	    success: boolean;
	    data?: academics.Section[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_academics_Section_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], academics.Section);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_admissions_Application_ {
	    success: boolean;
	    data?: admissions.Application[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_admissions_Application_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], admissions.Application);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_admissions_FormDefinition_ {
	    success: boolean;
	    data?: admissions.FormDefinition[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_admissions_FormDefinition_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], admissions.FormDefinition);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_admissions_FormVersion_ {
	    success: boolean;
	    data?: admissions.FormVersion[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_admissions_FormVersion_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], admissions.FormVersion);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_documents_Document_ {
	    success: boolean;
	    data?: documents.Document[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_documents_Document_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], documents.Document);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_fees_DiscountPolicy_ {
	    success: boolean;
	    data?: fees.DiscountPolicy[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_fees_DiscountPolicy_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], fees.DiscountPolicy);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_fees_FeeHead_ {
	    success: boolean;
	    data?: fees.FeeHead[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_fees_FeeHead_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], fees.FeeHead);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_fees_FeeStructureVersion_ {
	    success: boolean;
	    data?: fees.FeeStructureVersion[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_fees_FeeStructureVersion_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], fees.FeeStructureVersion);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_fees_FeeStructure_ {
	    success: boolean;
	    data?: fees.FeeStructure[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_fees_FeeStructure_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], fees.FeeStructure);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_fees_Invoice_ {
	    success: boolean;
	    data?: fees.Invoice[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_fees_Invoice_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], fees.Invoice);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_fees_StudentDiscount_ {
	    success: boolean;
	    data?: fees.StudentDiscount[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_fees_StudentDiscount_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], fees.StudentDiscount);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_guardians_StudentGuardianLink_ {
	    success: boolean;
	    data?: guardians.StudentGuardianLink[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_guardians_StudentGuardianLink_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], guardians.StudentGuardianLink);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_shared_ID_ {
	    success: boolean;
	    data?: string[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_shared_ID_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_students_Enrollment_ {
	    success: boolean;
	    data?: students.Enrollment[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_students_Enrollment_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], students.Enrollment);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_students_StatusHistoryEntry_ {
	    success: boolean;
	    data?: students.StatusHistoryEntry[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_students_StatusHistoryEntry_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], students.StatusHistoryEntry);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___school_management_internal_domain_students_Student_ {
	    success: boolean;
	    data?: students.Student[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response___school_management_internal_domain_students_Student_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = this.convertValues(source["data"], students.Student);
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response_bool_ {
	    success: boolean;
	    data?: boolean;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response_bool_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response_school_management_internal_domain_shared_ID_ {
	    success: boolean;
	    data?: string;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response_school_management_internal_domain_shared_ID_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response_string_ {
	    success: boolean;
	    data?: string;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Response_string_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = this.convertValues(source["error"], ErrorDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace enrollments {
	
	export class EnrollStudentInput {
	    studentId: string;
	    academicYearId: string;
	    classId: string;
	    sectionId?: string;
	    rollNo?: string;
	    enrolledOn: shared.Date;
	
	    static createFrom(source: any = {}) {
	        return new EnrollStudentInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.studentId = source["studentId"];
	        this.academicYearId = source["academicYearId"];
	        this.classId = source["classId"];
	        this.sectionId = source["sectionId"];
	        this.rollNo = source["rollNo"];
	        this.enrolledOn = this.convertValues(source["enrolledOn"], shared.Date);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PromotionItem {
	    studentId: string;
	    targetClassId: string;
	    targetSectionId?: string;
	    rollNo?: string;
	    action: string;
	
	    static createFrom(source: any = {}) {
	        return new PromotionItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.studentId = source["studentId"];
	        this.targetClassId = source["targetClassId"];
	        this.targetSectionId = source["targetSectionId"];
	        this.rollNo = source["rollNo"];
	        this.action = source["action"];
	    }
	}
	export class ExecutePromotionBatchInput {
	    fromAcademicYearId: string;
	    toAcademicYearId: string;
	    effectiveDate: shared.Date;
	    notes: string;
	    items: PromotionItem[];
	
	    static createFrom(source: any = {}) {
	        return new ExecutePromotionBatchInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fromAcademicYearId = source["fromAcademicYearId"];
	        this.toAcademicYearId = source["toAcademicYearId"];
	        this.effectiveDate = this.convertValues(source["effectiveDate"], shared.Date);
	        this.notes = source["notes"];
	        this.items = this.convertValues(source["items"], PromotionItem);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace fees {
	
	export class ApplyStudentDiscountInput {
	    studentId: string;
	    discountPolicyId: string;
	    academicYearId: string;
	    validFrom: shared.Date;
	    validTo?: shared.Date;
	
	    static createFrom(source: any = {}) {
	        return new ApplyStudentDiscountInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.studentId = source["studentId"];
	        this.discountPolicyId = source["discountPolicyId"];
	        this.academicYearId = source["academicYearId"];
	        this.validFrom = this.convertValues(source["validFrom"], shared.Date);
	        this.validTo = this.convertValues(source["validTo"], shared.Date);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CreateDiscountPolicyInput {
	    code: string;
	    name: string;
	    description: string;
	    discountType: string;
	    rateBasisPoints: number;
	    amountMinor: number;
	    feeHeadId?: string;
	    requiresApproval: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CreateDiscountPolicyInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.discountType = source["discountType"];
	        this.rateBasisPoints = source["rateBasisPoints"];
	        this.amountMinor = source["amountMinor"];
	        this.feeHeadId = source["feeHeadId"];
	        this.requiresApproval = source["requiresApproval"];
	    }
	}
	export class CreateFeeHeadInput {
	    code: string;
	    name: string;
	    description: string;
	    isRecurring: boolean;
	    defaultAmountMinor: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateFeeHeadInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.isRecurring = source["isRecurring"];
	        this.defaultAmountMinor = source["defaultAmountMinor"];
	    }
	}
	export class CreateFeeStructureInput {
	    academicYearId: string;
	    classId: string;
	    campusId?: string;
	    code: string;
	    name: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateFeeStructureInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.academicYearId = source["academicYearId"];
	        this.classId = source["classId"];
	        this.campusId = source["campusId"];
	        this.code = source["code"];
	        this.name = source["name"];
	        this.description = source["description"];
	    }
	}
	export class FeeStructureItemInput {
	    feeHeadId: string;
	    amountMinor: number;
	    frequency: string;
	    dueDay: number;
	    isOptional: boolean;
	
	    static createFrom(source: any = {}) {
	        return new FeeStructureItemInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.feeHeadId = source["feeHeadId"];
	        this.amountMinor = source["amountMinor"];
	        this.frequency = source["frequency"];
	        this.dueDay = source["dueDay"];
	        this.isOptional = source["isOptional"];
	    }
	}
	export class CreateFeeStructureVersionInput {
	    feeStructureId: string;
	    changelog: string;
	    items: FeeStructureItemInput[];
	
	    static createFrom(source: any = {}) {
	        return new CreateFeeStructureVersionInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.feeStructureId = source["feeStructureId"];
	        this.changelog = source["changelog"];
	        this.items = this.convertValues(source["items"], FeeStructureItemInput);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DecideStudentDiscountInput {
	    discountId: string;
	    approved: boolean;
	    rejectionReason: string;
	
	    static createFrom(source: any = {}) {
	        return new DecideStudentDiscountInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.discountId = source["discountId"];
	        this.approved = source["approved"];
	        this.rejectionReason = source["rejectionReason"];
	    }
	}
	export class DiscountPolicy {
	    ID: string;
	    SchoolID: string;
	    Code: string;
	    Name: string;
	    Description: string;
	    DiscountType: string;
	    RateBasisPoints: number;
	    AmountMinor: number;
	    FeeHeadID?: string;
	    RequiresApproval: boolean;
	    IsActive: boolean;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	
	    static createFrom(source: any = {}) {
	        return new DiscountPolicy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.SchoolID = source["SchoolID"];
	        this.Code = source["Code"];
	        this.Name = source["Name"];
	        this.Description = source["Description"];
	        this.DiscountType = source["DiscountType"];
	        this.RateBasisPoints = source["RateBasisPoints"];
	        this.AmountMinor = source["AmountMinor"];
	        this.FeeHeadID = source["FeeHeadID"];
	        this.RequiresApproval = source["RequiresApproval"];
	        this.IsActive = source["IsActive"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FeeHead {
	    ID: string;
	    SchoolID: string;
	    Code: string;
	    Name: string;
	    Description: string;
	    IsRecurring: boolean;
	    DefaultAmountMinor: number;
	    IsActive: boolean;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	
	    static createFrom(source: any = {}) {
	        return new FeeHead(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.SchoolID = source["SchoolID"];
	        this.Code = source["Code"];
	        this.Name = source["Name"];
	        this.Description = source["Description"];
	        this.IsRecurring = source["IsRecurring"];
	        this.DefaultAmountMinor = source["DefaultAmountMinor"];
	        this.IsActive = source["IsActive"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FeeStructure {
	    ID: string;
	    SchoolID: string;
	    AcademicYearID: string;
	    ClassID: string;
	    CampusID?: string;
	    Code: string;
	    Name: string;
	    Description: string;
	    IsActive: boolean;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	
	    static createFrom(source: any = {}) {
	        return new FeeStructure(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.SchoolID = source["SchoolID"];
	        this.AcademicYearID = source["AcademicYearID"];
	        this.ClassID = source["ClassID"];
	        this.CampusID = source["CampusID"];
	        this.Code = source["Code"];
	        this.Name = source["Name"];
	        this.Description = source["Description"];
	        this.IsActive = source["IsActive"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FeeStructureItem {
	    ID: string;
	    FeeStructureVersionID: string;
	    FeeHeadID: string;
	    AmountMinor: number;
	    Frequency: string;
	    DueDay: number;
	    IsOptional: boolean;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    FeeHeadCode: string;
	    FeeHeadName: string;
	
	    static createFrom(source: any = {}) {
	        return new FeeStructureItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.FeeStructureVersionID = source["FeeStructureVersionID"];
	        this.FeeHeadID = source["FeeHeadID"];
	        this.AmountMinor = source["AmountMinor"];
	        this.Frequency = source["Frequency"];
	        this.DueDay = source["DueDay"];
	        this.IsOptional = source["IsOptional"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.FeeHeadCode = source["FeeHeadCode"];
	        this.FeeHeadName = source["FeeHeadName"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class FeeStructureVersion {
	    ID: string;
	    FeeStructureID: string;
	    VersionNo: number;
	    Status: string;
	    Changelog: string;
	    CreatedByID: string;
	    PublishedByID?: string;
	    // Go type: time
	    PublishedAt?: any;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	
	    static createFrom(source: any = {}) {
	        return new FeeStructureVersion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.FeeStructureID = source["FeeStructureID"];
	        this.VersionNo = source["VersionNo"];
	        this.Status = source["Status"];
	        this.Changelog = source["Changelog"];
	        this.CreatedByID = source["CreatedByID"];
	        this.PublishedByID = source["PublishedByID"];
	        this.PublishedAt = this.convertValues(source["PublishedAt"], null);
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FeeStructureVersionDetail {
	    version: FeeStructureVersion;
	    items: FeeStructureItem[];
	
	    static createFrom(source: any = {}) {
	        return new FeeStructureVersionDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = this.convertValues(source["version"], FeeStructureVersion);
	        this.items = this.convertValues(source["items"], FeeStructureItem);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GenerateInvoicesBatchInput {
	    academicYearId: string;
	    classId?: string;
	    billingPeriod: string;
	    issueDate: shared.Date;
	    dueDate: shared.Date;
	
	    static createFrom(source: any = {}) {
	        return new GenerateInvoicesBatchInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.academicYearId = source["academicYearId"];
	        this.classId = source["classId"];
	        this.billingPeriod = source["billingPeriod"];
	        this.issueDate = this.convertValues(source["issueDate"], shared.Date);
	        this.dueDate = this.convertValues(source["dueDate"], shared.Date);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Invoice {
	    ID: string;
	    SchoolID: string;
	    StudentID: string;
	    AcademicYearID: string;
	    InvoiceNo: string;
	    IssueDate: shared.Date;
	    DueDate: shared.Date;
	    BillingPeriod: string;
	    SubtotalAmountMinor: number;
	    DiscountAmountMinor: number;
	    TotalAmountMinor: number;
	    PaidAmountMinor: number;
	    BalanceMinor: number;
	    Status: string;
	    VoidReason: string;
	    VoidedByID?: string;
	    // Go type: time
	    VoidedAt?: any;
	    CreatedByID: string;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	    StudentFirstName: string;
	    StudentLastName: string;
	    AdmissionNo: string;
	
	    static createFrom(source: any = {}) {
	        return new Invoice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.SchoolID = source["SchoolID"];
	        this.StudentID = source["StudentID"];
	        this.AcademicYearID = source["AcademicYearID"];
	        this.InvoiceNo = source["InvoiceNo"];
	        this.IssueDate = this.convertValues(source["IssueDate"], shared.Date);
	        this.DueDate = this.convertValues(source["DueDate"], shared.Date);
	        this.BillingPeriod = source["BillingPeriod"];
	        this.SubtotalAmountMinor = source["SubtotalAmountMinor"];
	        this.DiscountAmountMinor = source["DiscountAmountMinor"];
	        this.TotalAmountMinor = source["TotalAmountMinor"];
	        this.PaidAmountMinor = source["PaidAmountMinor"];
	        this.BalanceMinor = source["BalanceMinor"];
	        this.Status = source["Status"];
	        this.VoidReason = source["VoidReason"];
	        this.VoidedByID = source["VoidedByID"];
	        this.VoidedAt = this.convertValues(source["VoidedAt"], null);
	        this.CreatedByID = source["CreatedByID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	        this.StudentFirstName = source["StudentFirstName"];
	        this.StudentLastName = source["StudentLastName"];
	        this.AdmissionNo = source["AdmissionNo"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class InvoiceItem {
	    ID: string;
	    InvoiceID: string;
	    FeeHeadID: string;
	    FeeHeadName: string;
	    Description: string;
	    OriginalAmountMinor: number;
	    DiscountAmountMinor: number;
	    NetAmountMinor: number;
	    // Go type: time
	    CreatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new InvoiceItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.InvoiceID = source["InvoiceID"];
	        this.FeeHeadID = source["FeeHeadID"];
	        this.FeeHeadName = source["FeeHeadName"];
	        this.Description = source["Description"];
	        this.OriginalAmountMinor = source["OriginalAmountMinor"];
	        this.DiscountAmountMinor = source["DiscountAmountMinor"];
	        this.NetAmountMinor = source["NetAmountMinor"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class InvoiceDetail {
	    invoice: Invoice;
	    items: InvoiceItem[];
	
	    static createFrom(source: any = {}) {
	        return new InvoiceDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.invoice = this.convertValues(source["invoice"], Invoice);
	        this.items = this.convertValues(source["items"], InvoiceItem);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class InvoiceFilter {
	    SchoolID: string;
	    Status: string;
	    AcademicYearID: string;
	    StudentID: string;
	    Limit: number;
	
	    static createFrom(source: any = {}) {
	        return new InvoiceFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SchoolID = source["SchoolID"];
	        this.Status = source["Status"];
	        this.AcademicYearID = source["AcademicYearID"];
	        this.StudentID = source["StudentID"];
	        this.Limit = source["Limit"];
	    }
	}
	
	export class StudentDiscount {
	    ID: string;
	    StudentID: string;
	    DiscountPolicyID: string;
	    AcademicYearID: string;
	    Status: string;
	    ApprovedByID?: string;
	    // Go type: time
	    ApprovedAt?: any;
	    RejectionReason: string;
	    ValidFrom: shared.Date;
	    ValidTo?: shared.Date;
	    CreatedByID: string;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	    PolicyCode: string;
	    PolicyName: string;
	    DiscountType: string;
	    RateBasisPoints: number;
	    PolicyAmountMinor: number;
	    PolicyFeeHeadID?: string;
	    StudentFirstName: string;
	    StudentLastName: string;
	    AdmissionNo: string;
	
	    static createFrom(source: any = {}) {
	        return new StudentDiscount(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.StudentID = source["StudentID"];
	        this.DiscountPolicyID = source["DiscountPolicyID"];
	        this.AcademicYearID = source["AcademicYearID"];
	        this.Status = source["Status"];
	        this.ApprovedByID = source["ApprovedByID"];
	        this.ApprovedAt = this.convertValues(source["ApprovedAt"], null);
	        this.RejectionReason = source["RejectionReason"];
	        this.ValidFrom = this.convertValues(source["ValidFrom"], shared.Date);
	        this.ValidTo = this.convertValues(source["ValidTo"], shared.Date);
	        this.CreatedByID = source["CreatedByID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	        this.PolicyCode = source["PolicyCode"];
	        this.PolicyName = source["PolicyName"];
	        this.DiscountType = source["DiscountType"];
	        this.RateBasisPoints = source["RateBasisPoints"];
	        this.PolicyAmountMinor = source["PolicyAmountMinor"];
	        this.PolicyFeeHeadID = source["PolicyFeeHeadID"];
	        this.StudentFirstName = source["StudentFirstName"];
	        this.StudentLastName = source["StudentLastName"];
	        this.AdmissionNo = source["AdmissionNo"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateFeeHeadInput {
	    id: string;
	    name: string;
	    description: string;
	    isRecurring: boolean;
	    defaultAmountMinor: number;
	    isActive: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UpdateFeeHeadInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.isRecurring = source["isRecurring"];
	        this.defaultAmountMinor = source["defaultAmountMinor"];
	        this.isActive = source["isActive"];
	    }
	}
	export class VoidInvoiceInput {
	    invoiceId: string;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new VoidInvoiceInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.invoiceId = source["invoiceId"];
	        this.reason = source["reason"];
	    }
	}

}

export namespace guardians {
	
	export class CreateGuardianInput {
	    firstName: string;
	    lastName: string;
	    phonePrimary: string;
	    phoneSecondary: string;
	    email: string;
	    occupation: string;
	    nationalId: string;
	    address: shared.Address;
	
	    static createFrom(source: any = {}) {
	        return new CreateGuardianInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.firstName = source["firstName"];
	        this.lastName = source["lastName"];
	        this.phonePrimary = source["phonePrimary"];
	        this.phoneSecondary = source["phoneSecondary"];
	        this.email = source["email"];
	        this.occupation = source["occupation"];
	        this.nationalId = source["nationalId"];
	        this.address = this.convertValues(source["address"], shared.Address);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Guardian {
	    ID: string;
	    SchoolID: string;
	    FirstName: string;
	    LastName: string;
	    PhonePrimary: string;
	    PhonePrimaryNormalized: string;
	    PhoneSecondary: string;
	    Email: string;
	    Occupation: string;
	    NationalID: string;
	    Address: shared.Address;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	
	    static createFrom(source: any = {}) {
	        return new Guardian(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.SchoolID = source["SchoolID"];
	        this.FirstName = source["FirstName"];
	        this.LastName = source["LastName"];
	        this.PhonePrimary = source["PhonePrimary"];
	        this.PhonePrimaryNormalized = source["PhonePrimaryNormalized"];
	        this.PhoneSecondary = source["PhoneSecondary"];
	        this.Email = source["Email"];
	        this.Occupation = source["Occupation"];
	        this.NationalID = source["NationalID"];
	        this.Address = this.convertValues(source["Address"], shared.Address);
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StudentGuardianLinkInput {
	    guardianId: string;
	    relationship: string;
	    isPrimary: boolean;
	    isEmergencyContact: boolean;
	    receivesBilling: boolean;
	    canPickup: boolean;
	    priority: number;
	
	    static createFrom(source: any = {}) {
	        return new StudentGuardianLinkInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.guardianId = source["guardianId"];
	        this.relationship = source["relationship"];
	        this.isPrimary = source["isPrimary"];
	        this.isEmergencyContact = source["isEmergencyContact"];
	        this.receivesBilling = source["receivesBilling"];
	        this.canPickup = source["canPickup"];
	        this.priority = source["priority"];
	    }
	}
	export class SetStudentGuardiansInput {
	    studentId: string;
	    guardians: StudentGuardianLinkInput[];
	
	    static createFrom(source: any = {}) {
	        return new SetStudentGuardiansInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.studentId = source["studentId"];
	        this.guardians = this.convertValues(source["guardians"], StudentGuardianLinkInput);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StudentGuardianLink {
	    StudentID: string;
	    GuardianID: string;
	    Relationship: string;
	    IsPrimary: boolean;
	    IsEmergencyContact: boolean;
	    ReceivesBilling: boolean;
	    CanPickup: boolean;
	    Priority: number;
	    FirstName: string;
	    LastName: string;
	    PhonePrimary: string;
	    PhoneSecondary: string;
	    Email: string;
	    Occupation: string;
	
	    static createFrom(source: any = {}) {
	        return new StudentGuardianLink(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.StudentID = source["StudentID"];
	        this.GuardianID = source["GuardianID"];
	        this.Relationship = source["Relationship"];
	        this.IsPrimary = source["IsPrimary"];
	        this.IsEmergencyContact = source["IsEmergencyContact"];
	        this.ReceivesBilling = source["ReceivesBilling"];
	        this.CanPickup = source["CanPickup"];
	        this.Priority = source["Priority"];
	        this.FirstName = source["FirstName"];
	        this.LastName = source["LastName"];
	        this.PhonePrimary = source["PhonePrimary"];
	        this.PhoneSecondary = source["PhoneSecondary"];
	        this.Email = source["Email"];
	        this.Occupation = source["Occupation"];
	    }
	}
	
	export class UpdateGuardianInput {
	    id: string;
	    firstName: string;
	    lastName: string;
	    phonePrimary: string;
	    phoneSecondary: string;
	    email: string;
	    occupation: string;
	    nationalId: string;
	    address: shared.Address;
	
	    static createFrom(source: any = {}) {
	        return new UpdateGuardianInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.firstName = source["firstName"];
	        this.lastName = source["lastName"];
	        this.phonePrimary = source["phonePrimary"];
	        this.phoneSecondary = source["phoneSecondary"];
	        this.email = source["email"];
	        this.occupation = source["occupation"];
	        this.nationalId = source["nationalId"];
	        this.address = this.convertValues(source["address"], shared.Address);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace payments {
	
	export class AdjustmentFilter {
	    SchoolID: string;
	    StudentID?: string;
	    Status?: string;
	    Limit: number;
	
	    static createFrom(source: any = {}) {
	        return new AdjustmentFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SchoolID = source["SchoolID"];
	        this.StudentID = source["StudentID"];
	        this.Status = source["Status"];
	        this.Limit = source["Limit"];
	    }
	}
	export class CreditLedgerEntry {
	    ID: string;
	    CreditBalanceID: string;
	    StudentID: string;
	    AmountMinor: number;
	    BalanceAfterMinor: number;
	    ReferenceType: string;
	    ReferenceID?: string;
	    Notes: string;
	    // Go type: time
	    CreatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new CreditLedgerEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreditBalanceID = source["CreditBalanceID"];
	        this.StudentID = source["StudentID"];
	        this.AmountMinor = source["AmountMinor"];
	        this.BalanceAfterMinor = source["BalanceAfterMinor"];
	        this.ReferenceType = source["ReferenceType"];
	        this.ReferenceID = source["ReferenceID"];
	        this.Notes = source["Notes"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DecideFinancialAdjustmentInput {
	    adjustmentId: string;
	    approved: boolean;
	    decisionNotes: string;
	
	    static createFrom(source: any = {}) {
	        return new DecideFinancialAdjustmentInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.adjustmentId = source["adjustmentId"];
	        this.approved = source["approved"];
	        this.decisionNotes = source["decisionNotes"];
	    }
	}
	export class FinancialAdjustment {
	    ID: string;
	    SchoolID: string;
	    StudentID: string;
	    InvoiceID?: string;
	    InvoiceItemID?: string;
	    AdjustmentType: string;
	    AmountMinor: number;
	    Reason: string;
	    Status: string;
	    RequestedByID: string;
	    // Go type: time
	    RequestedAt: any;
	    DecidedByID?: string;
	    // Go type: time
	    DecidedAt?: any;
	    DecisionNotes: string;
	    StudentFirstName: string;
	    StudentLastName: string;
	    AdmissionNo: string;
	    RequestedByName: string;
	    DecidedByName?: string;
	    InvoiceNo?: string;
	
	    static createFrom(source: any = {}) {
	        return new FinancialAdjustment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.SchoolID = source["SchoolID"];
	        this.StudentID = source["StudentID"];
	        this.InvoiceID = source["InvoiceID"];
	        this.InvoiceItemID = source["InvoiceItemID"];
	        this.AdjustmentType = source["AdjustmentType"];
	        this.AmountMinor = source["AmountMinor"];
	        this.Reason = source["Reason"];
	        this.Status = source["Status"];
	        this.RequestedByID = source["RequestedByID"];
	        this.RequestedAt = this.convertValues(source["RequestedAt"], null);
	        this.DecidedByID = source["DecidedByID"];
	        this.DecidedAt = this.convertValues(source["DecidedAt"], null);
	        this.DecisionNotes = source["DecisionNotes"];
	        this.StudentFirstName = source["StudentFirstName"];
	        this.StudentLastName = source["StudentLastName"];
	        this.AdmissionNo = source["AdmissionNo"];
	        this.RequestedByName = source["RequestedByName"];
	        this.DecidedByName = source["DecidedByName"];
	        this.InvoiceNo = source["InvoiceNo"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PaymentAllocation {
	    ID: string;
	    PaymentID: string;
	    InvoiceID: string;
	    InvoiceItemID?: string;
	    AllocatedAmountMinor: number;
	    // Go type: time
	    AllocatedAt: any;
	    AllocatedByID: string;
	    InvoiceNo: string;
	    BillingPeriod: string;
	    ReceiptNo: string;
	    PaymentDate: shared.Date;
	    PaymentMethod: string;
	
	    static createFrom(source: any = {}) {
	        return new PaymentAllocation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.PaymentID = source["PaymentID"];
	        this.InvoiceID = source["InvoiceID"];
	        this.InvoiceItemID = source["InvoiceItemID"];
	        this.AllocatedAmountMinor = source["AllocatedAmountMinor"];
	        this.AllocatedAt = this.convertValues(source["AllocatedAt"], null);
	        this.AllocatedByID = source["AllocatedByID"];
	        this.InvoiceNo = source["InvoiceNo"];
	        this.BillingPeriod = source["BillingPeriod"];
	        this.ReceiptNo = source["ReceiptNo"];
	        this.PaymentDate = this.convertValues(source["PaymentDate"], shared.Date);
	        this.PaymentMethod = source["PaymentMethod"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PaymentReversal {
	    ID: string;
	    PaymentID: string;
	    ReversalDate: shared.Date;
	    Reason: string;
	    ReversedByID: string;
	    // Go type: time
	    ReversedAt: any;
	    ReversedByName: string;
	
	    static createFrom(source: any = {}) {
	        return new PaymentReversal(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.PaymentID = source["PaymentID"];
	        this.ReversalDate = this.convertValues(source["ReversalDate"], shared.Date);
	        this.Reason = source["Reason"];
	        this.ReversedByID = source["ReversedByID"];
	        this.ReversedAt = this.convertValues(source["ReversedAt"], null);
	        this.ReversedByName = source["ReversedByName"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PaymentWithDetails {
	    ID: string;
	    SchoolID: string;
	    StudentID: string;
	    ReceiptNo: string;
	    PaymentDate: shared.Date;
	    PaymentMethod: string;
	    AmountMinor: number;
	    AllocatedAmountMinor: number;
	    UnallocatedAmountMinor: number;
	    ReferenceNo: string;
	    PayerName: string;
	    PayerPhone: string;
	    Status: string;
	    RecordedByID: string;
	    // Go type: time
	    RecordedAt: any;
	    Notes: string;
	    StudentFirstName: string;
	    StudentLastName: string;
	    AdmissionNo: string;
	    RecordedByName: string;
	
	    static createFrom(source: any = {}) {
	        return new PaymentWithDetails(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.SchoolID = source["SchoolID"];
	        this.StudentID = source["StudentID"];
	        this.ReceiptNo = source["ReceiptNo"];
	        this.PaymentDate = this.convertValues(source["PaymentDate"], shared.Date);
	        this.PaymentMethod = source["PaymentMethod"];
	        this.AmountMinor = source["AmountMinor"];
	        this.AllocatedAmountMinor = source["AllocatedAmountMinor"];
	        this.UnallocatedAmountMinor = source["UnallocatedAmountMinor"];
	        this.ReferenceNo = source["ReferenceNo"];
	        this.PayerName = source["PayerName"];
	        this.PayerPhone = source["PayerPhone"];
	        this.Status = source["Status"];
	        this.RecordedByID = source["RecordedByID"];
	        this.RecordedAt = this.convertValues(source["RecordedAt"], null);
	        this.Notes = source["Notes"];
	        this.StudentFirstName = source["StudentFirstName"];
	        this.StudentLastName = source["StudentLastName"];
	        this.AdmissionNo = source["AdmissionNo"];
	        this.RecordedByName = source["RecordedByName"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PaymentDetail {
	    payment: PaymentWithDetails;
	    allocations: PaymentAllocation[];
	    reversal?: PaymentReversal;
	
	    static createFrom(source: any = {}) {
	        return new PaymentDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.payment = this.convertValues(source["payment"], PaymentWithDetails);
	        this.allocations = this.convertValues(source["allocations"], PaymentAllocation);
	        this.reversal = this.convertValues(source["reversal"], PaymentReversal);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PaymentFilter {
	    SchoolID: string;
	    StudentID?: string;
	    Status?: string;
	    PaymentMethod?: string;
	    Limit: number;
	
	    static createFrom(source: any = {}) {
	        return new PaymentFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SchoolID = source["SchoolID"];
	        this.StudentID = source["StudentID"];
	        this.Status = source["Status"];
	        this.PaymentMethod = source["PaymentMethod"];
	        this.Limit = source["Limit"];
	    }
	}
	
	
	export class RecordPaymentInput {
	    studentId: string;
	    paymentDate: string;
	    paymentMethod: string;
	    amountMinor: number;
	    referenceNo: string;
	    payerName: string;
	    payerPhone: string;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new RecordPaymentInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.studentId = source["studentId"];
	        this.paymentDate = source["paymentDate"];
	        this.paymentMethod = source["paymentMethod"];
	        this.amountMinor = source["amountMinor"];
	        this.referenceNo = source["referenceNo"];
	        this.payerName = source["payerName"];
	        this.payerPhone = source["payerPhone"];
	        this.notes = source["notes"];
	    }
	}
	export class RecordPaymentResult {
	    paymentId: string;
	    receiptNo: string;
	    allocatedMinor: number;
	    unallocatedMinor: number;
	
	    static createFrom(source: any = {}) {
	        return new RecordPaymentResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.paymentId = source["paymentId"];
	        this.receiptNo = source["receiptNo"];
	        this.allocatedMinor = source["allocatedMinor"];
	        this.unallocatedMinor = source["unallocatedMinor"];
	    }
	}
	export class RequestFinancialAdjustmentInput {
	    studentId: string;
	    invoiceId?: string;
	    invoiceItemId?: string;
	    adjustmentType: string;
	    amountMinor: number;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new RequestFinancialAdjustmentInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.studentId = source["studentId"];
	        this.invoiceId = source["invoiceId"];
	        this.invoiceItemId = source["invoiceItemId"];
	        this.adjustmentType = source["adjustmentType"];
	        this.amountMinor = source["amountMinor"];
	        this.reason = source["reason"];
	    }
	}
	export class ReversePaymentInput {
	    paymentId: string;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new ReversePaymentInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.paymentId = source["paymentId"];
	        this.reason = source["reason"];
	    }
	}
	export class StudentCreditBalanceInfo {
	    balanceMinor: number;
	    ledger: CreditLedgerEntry[];
	
	    static createFrom(source: any = {}) {
	        return new StudentCreditBalanceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.balanceMinor = source["balanceMinor"];
	        this.ledger = this.convertValues(source["ledger"], CreditLedgerEntry);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace roles {
	
	export class CreateRoleInput {
	    code: string;
	    name: string;
	    description: string;
	    permissions: string[];
	
	    static createFrom(source: any = {}) {
	        return new CreateRoleInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.permissions = source["permissions"];
	    }
	}
	export class PermissionInfoDTO {
	    code: string;
	    module: string;
	    description: string;
	    sensitive: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PermissionInfoDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.module = source["module"];
	        this.description = source["description"];
	        this.sensitive = source["sensitive"];
	    }
	}
	export class RoleDTO {
	    id: string;
	    code: string;
	    name: string;
	    description: string;
	    isSystem: boolean;
	    permissions: string[];
	    userCount: number;
	
	    static createFrom(source: any = {}) {
	        return new RoleDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.code = source["code"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.isSystem = source["isSystem"];
	        this.permissions = source["permissions"];
	        this.userCount = source["userCount"];
	    }
	}
	export class UpdateRoleInput {
	    id: string;
	    name: string;
	    description: string;
	    permissions: string[];
	
	    static createFrom(source: any = {}) {
	        return new UpdateRoleInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.permissions = source["permissions"];
	    }
	}

}

export namespace school {
	
	export class School {
	    ID: string;
	    Name: string;
	    LegalName: string;
	    Code: string;
	    LogoDocumentID?: string;
	    Address: shared.Address;
	    Phone: string;
	    Email: string;
	    Website: string;
	    Timezone: string;
	    CurrencyCode: string;
	    CurrencyMinorUnits: number;
	    Locale: string;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	
	    static createFrom(source: any = {}) {
	        return new School(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.LegalName = source["LegalName"];
	        this.Code = source["Code"];
	        this.LogoDocumentID = source["LogoDocumentID"];
	        this.Address = this.convertValues(source["Address"], shared.Address);
	        this.Phone = source["Phone"];
	        this.Email = source["Email"];
	        this.Website = source["Website"];
	        this.Timezone = source["Timezone"];
	        this.CurrencyCode = source["CurrencyCode"];
	        this.CurrencyMinorUnits = source["CurrencyMinorUnits"];
	        this.Locale = source["Locale"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Settings {
	    SchoolID: string;
	    AdmissionNoFormat: string;
	    InvoicePrefix: string;
	    ReceiptPrefix: string;
	    RoundingUnitMinor: number;
	    RoundingMode: string;
	    DefaultDueDay: number;
	    LateFeeGraceDays: number;
	    DiscountCombination: string;
	    RequireSecondApprover: boolean;
	    DateFormat: string;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SchoolID = source["SchoolID"];
	        this.AdmissionNoFormat = source["AdmissionNoFormat"];
	        this.InvoicePrefix = source["InvoicePrefix"];
	        this.ReceiptPrefix = source["ReceiptPrefix"];
	        this.RoundingUnitMinor = source["RoundingUnitMinor"];
	        this.RoundingMode = source["RoundingMode"];
	        this.DefaultDueDay = source["DefaultDueDay"];
	        this.LateFeeGraceDays = source["LateFeeGraceDays"];
	        this.DiscountCombination = source["DiscountCombination"];
	        this.RequireSecondApprover = source["RequireSecondApprover"];
	        this.DateFormat = source["DateFormat"];
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateSchoolProfileInput {
	    name: string;
	    legalName: string;
	    code: string;
	    address: shared.Address;
	    phone: string;
	    email: string;
	    website: string;
	    timezone: string;
	    locale: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateSchoolProfileInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.legalName = source["legalName"];
	        this.code = source["code"];
	        this.address = this.convertValues(source["address"], shared.Address);
	        this.phone = source["phone"];
	        this.email = source["email"];
	        this.website = source["website"];
	        this.timezone = source["timezone"];
	        this.locale = source["locale"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateSchoolSettingsInput {
	    admissionNoFormat: string;
	    invoicePrefix: string;
	    receiptPrefix: string;
	    roundingUnitMinor: number;
	    roundingMode: string;
	    defaultDueDay: number;
	    lateFeeGraceDays: number;
	    discountCombination: string;
	    requireSecondApprover: boolean;
	    dateFormat: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateSchoolSettingsInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.admissionNoFormat = source["admissionNoFormat"];
	        this.invoicePrefix = source["invoicePrefix"];
	        this.receiptPrefix = source["receiptPrefix"];
	        this.roundingUnitMinor = source["roundingUnitMinor"];
	        this.roundingMode = source["roundingMode"];
	        this.defaultDueDay = source["defaultDueDay"];
	        this.lateFeeGraceDays = source["lateFeeGraceDays"];
	        this.discountCombination = source["discountCombination"];
	        this.requireSecondApprover = source["requireSecondApprover"];
	        this.dateFormat = source["dateFormat"];
	    }
	}

}

export namespace shared {
	
	export class Address {
	    line1?: string;
	    line2?: string;
	    city?: string;
	    region?: string;
	    postalCode?: string;
	    country?: string;
	
	    static createFrom(source: any = {}) {
	        return new Address(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.line1 = source["line1"];
	        this.line2 = source["line2"];
	        this.city = source["city"];
	        this.region = source["region"];
	        this.postalCode = source["postalCode"];
	        this.country = source["country"];
	    }
	}
	export class Date {
	    Year: number;
	    Month: number;
	    Day: number;
	
	    static createFrom(source: any = {}) {
	        return new Date(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Year = source["Year"];
	        this.Month = source["Month"];
	        this.Day = source["Day"];
	    }
	}

}

export namespace students {
	
	export class ChangeStudentStatusInput {
	    studentId: string;
	    newStatus: string;
	    effectiveOn: shared.Date;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new ChangeStudentStatusInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.studentId = source["studentId"];
	        this.newStatus = source["newStatus"];
	        this.effectiveOn = this.convertValues(source["effectiveOn"], shared.Date);
	        this.reason = source["reason"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Enrollment {
	    ID: string;
	    StudentID: string;
	    AcademicYearID: string;
	    ClassID: string;
	    SectionID?: string;
	    RollNo?: string;
	    Status: string;
	    EnrolledOn: shared.Date;
	    EndedOn?: shared.Date;
	    PreviousEnrollmentID?: string;
	    PromotionBatchID?: string;
	    AdmissionApplicationID?: string;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	    ClassName: string;
	    SectionName: string;
	    YearCode: string;
	
	    static createFrom(source: any = {}) {
	        return new Enrollment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.StudentID = source["StudentID"];
	        this.AcademicYearID = source["AcademicYearID"];
	        this.ClassID = source["ClassID"];
	        this.SectionID = source["SectionID"];
	        this.RollNo = source["RollNo"];
	        this.Status = source["Status"];
	        this.EnrolledOn = this.convertValues(source["EnrolledOn"], shared.Date);
	        this.EndedOn = this.convertValues(source["EndedOn"], shared.Date);
	        this.PreviousEnrollmentID = source["PreviousEnrollmentID"];
	        this.PromotionBatchID = source["PromotionBatchID"];
	        this.AdmissionApplicationID = source["AdmissionApplicationID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	        this.ClassName = source["ClassName"];
	        this.SectionName = source["SectionName"];
	        this.YearCode = source["YearCode"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PromotionBatch {
	    ID: string;
	    FromAcademicYearID: string;
	    ToAcademicYearID: string;
	    Status: string;
	    SummaryJSON: string;
	    ExecutedBy?: string;
	    // Go type: time
	    ExecutedAt?: any;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	
	    static createFrom(source: any = {}) {
	        return new PromotionBatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.FromAcademicYearID = source["FromAcademicYearID"];
	        this.ToAcademicYearID = source["ToAcademicYearID"];
	        this.Status = source["Status"];
	        this.SummaryJSON = source["SummaryJSON"];
	        this.ExecutedBy = source["ExecutedBy"];
	        this.ExecutedAt = this.convertValues(source["ExecutedAt"], null);
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RegisterStudentInput {
	    admissionNo: string;
	    firstName: string;
	    middleName: string;
	    lastName: string;
	    dateOfBirth: shared.Date;
	    gender: string;
	    nationalId: string;
	    address: shared.Address;
	    photoDocId?: string;
	    admissionDate: shared.Date;
	    notes: string;
	    academicYearId: string;
	    classId: string;
	    sectionId?: string;
	    rollNo?: string;
	    admissionApplicationId?: string;
	
	    static createFrom(source: any = {}) {
	        return new RegisterStudentInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.admissionNo = source["admissionNo"];
	        this.firstName = source["firstName"];
	        this.middleName = source["middleName"];
	        this.lastName = source["lastName"];
	        this.dateOfBirth = this.convertValues(source["dateOfBirth"], shared.Date);
	        this.gender = source["gender"];
	        this.nationalId = source["nationalId"];
	        this.address = this.convertValues(source["address"], shared.Address);
	        this.photoDocId = source["photoDocId"];
	        this.admissionDate = this.convertValues(source["admissionDate"], shared.Date);
	        this.notes = source["notes"];
	        this.academicYearId = source["academicYearId"];
	        this.classId = source["classId"];
	        this.sectionId = source["sectionId"];
	        this.rollNo = source["rollNo"];
	        this.admissionApplicationId = source["admissionApplicationId"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StatusHistoryEntry {
	    ID: string;
	    StudentID: string;
	    FromStatus: string;
	    ToStatus: string;
	    EffectiveOn: shared.Date;
	    Reason: string;
	    ChangedBy: string;
	    // Go type: time
	    ChangedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new StatusHistoryEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.StudentID = source["StudentID"];
	        this.FromStatus = source["FromStatus"];
	        this.ToStatus = source["ToStatus"];
	        this.EffectiveOn = this.convertValues(source["EffectiveOn"], shared.Date);
	        this.Reason = source["Reason"];
	        this.ChangedBy = source["ChangedBy"];
	        this.ChangedAt = this.convertValues(source["ChangedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Student {
	    ID: string;
	    SchoolID: string;
	    AdmissionNo: string;
	    FirstName: string;
	    MiddleName: string;
	    LastName: string;
	    SearchName: string;
	    DateOfBirth: shared.Date;
	    Gender: string;
	    NationalID: string;
	    Address: shared.Address;
	    PhotoDocID?: string;
	    AdmissionDate: shared.Date;
	    Status: string;
	    Notes: string;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    Revision: number;
	
	    static createFrom(source: any = {}) {
	        return new Student(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.SchoolID = source["SchoolID"];
	        this.AdmissionNo = source["AdmissionNo"];
	        this.FirstName = source["FirstName"];
	        this.MiddleName = source["MiddleName"];
	        this.LastName = source["LastName"];
	        this.SearchName = source["SearchName"];
	        this.DateOfBirth = this.convertValues(source["DateOfBirth"], shared.Date);
	        this.Gender = source["Gender"];
	        this.NationalID = source["NationalID"];
	        this.Address = this.convertValues(source["Address"], shared.Address);
	        this.PhotoDocID = source["PhotoDocID"];
	        this.AdmissionDate = this.convertValues(source["AdmissionDate"], shared.Date);
	        this.Status = source["Status"];
	        this.Notes = source["Notes"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Revision = source["Revision"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StudentFilter {
	    SchoolID: string;
	    Status: string;
	    ClassID: string;
	    SectionID: string;
	    SearchQuery: string;
	    CursorName: string;
	    CursorID: string;
	    Limit: number;
	
	    static createFrom(source: any = {}) {
	        return new StudentFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SchoolID = source["SchoolID"];
	        this.Status = source["Status"];
	        this.ClassID = source["ClassID"];
	        this.SectionID = source["SectionID"];
	        this.SearchQuery = source["SearchQuery"];
	        this.CursorName = source["CursorName"];
	        this.CursorID = source["CursorID"];
	        this.Limit = source["Limit"];
	    }
	}
	export class UpdateStudentInput {
	    id: string;
	    firstName: string;
	    middleName: string;
	    lastName: string;
	    dateOfBirth: shared.Date;
	    gender: string;
	    nationalId: string;
	    address: shared.Address;
	    photoDocId?: string;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateStudentInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.firstName = source["firstName"];
	        this.middleName = source["middleName"];
	        this.lastName = source["lastName"];
	        this.dateOfBirth = this.convertValues(source["dateOfBirth"], shared.Date);
	        this.gender = source["gender"];
	        this.nationalId = source["nationalId"];
	        this.address = this.convertValues(source["address"], shared.Address);
	        this.photoDocId = source["photoDocId"];
	        this.notes = source["notes"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace users {
	
	export class AdminResetPasswordInput {
	    userId: string;
	    newPassword: string;
	
	    static createFrom(source: any = {}) {
	        return new AdminResetPasswordInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.userId = source["userId"];
	        this.newPassword = source["newPassword"];
	    }
	}
	export class CreateUserInput {
	    username: string;
	    displayName: string;
	    email: string;
	    password: string;
	    roleIds: string[];
	
	    static createFrom(source: any = {}) {
	        return new CreateUserInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.displayName = source["displayName"];
	        this.email = source["email"];
	        this.password = source["password"];
	        this.roleIds = source["roleIds"];
	    }
	}
	export class UpdateUserInput {
	    id: string;
	    displayName: string;
	    email: string;
	    roleIds: string[];
	
	    static createFrom(source: any = {}) {
	        return new UpdateUserInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.displayName = source["displayName"];
	        this.email = source["email"];
	        this.roleIds = source["roleIds"];
	    }
	}
	export class UserListItemDTO {
	    id: string;
	    username: string;
	    displayName: string;
	    email: string;
	    status: string;
	    roles: string[];
	    roleIds: string[];
	    mustChangePassword: boolean;
	    lastLoginAt?: string;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new UserListItemDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	        this.displayName = source["displayName"];
	        this.email = source["email"];
	        this.status = source["status"];
	        this.roles = source["roles"];
	        this.roleIds = source["roleIds"];
	        this.mustChangePassword = source["mustChangePassword"];
	        this.lastLoginAt = source["lastLoginAt"];
	        this.createdAt = source["createdAt"];
	    }
	}

}

export namespace wails {
	
	export class DocumentContentResponse {
	    document?: documents.Document;
	    dataBase64: string;
	
	    static createFrom(source: any = {}) {
	        return new DocumentContentResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.document = this.convertValues(source["document"], documents.Document);
	        this.dataBase64 = source["dataBase64"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UploadDocumentRequest {
	    studentId?: string;
	    applicationId?: string;
	    documentTypeCode: string;
	    fileName: string;
	    mimeType: string;
	    dataBase64: string;
	
	    static createFrom(source: any = {}) {
	        return new UploadDocumentRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.studentId = source["studentId"];
	        this.applicationId = source["applicationId"];
	        this.documentTypeCode = source["documentTypeCode"];
	        this.fileName = source["fileName"];
	        this.mimeType = source["mimeType"];
	        this.dataBase64 = source["dataBase64"];
	    }
	}

}

