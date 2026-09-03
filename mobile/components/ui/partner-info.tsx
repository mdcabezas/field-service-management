import React, { useState } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { Card } from './card';
import { Badge } from './badge';
import { Button } from './button';

interface Partner {
  id: string;
  name: string;
  order_id: string;
  supervisor: string;
  contacts?: PartnerContact[];
}

interface PartnerContact {
  id: string;
  name: string;
  role: string;
  phone: string;
  email: string;
}

interface PartnerInfoProps {
  partner: Partner;
  showContacts?: boolean;
}

export const PartnerInfo = React.memo(function PartnerInfo({ partner, showContacts = false }: PartnerInfoProps) {
  const [contactsExpanded, setContactsExpanded] = useState(false);

  return (
    <Card title="Partner">
      <View style={styles.info}>
        <View style={styles.row}>
          <Text style={styles.label}>Nombre:</Text>
          <Text style={styles.value}>{partner.name}</Text>
        </View>
        <View style={styles.row}>
          <Text style={styles.label}>Orden:</Text>
          <Text style={styles.value}>{partner.order_id}</Text>
        </View>
        <View style={styles.row}>
          <Text style={styles.label}>Supervisor:</Text>
          <Text style={styles.value}>{partner.supervisor}</Text>
        </View>
      </View>

      {showContacts && partner.contacts && partner.contacts.length > 0 && (
        <View style={styles.contactsSection}>
          <Button
            title={contactsExpanded ? 'Ocultar contactos' : `Ver contactos (${partner.contacts.length})`}
            onPress={() => setContactsExpanded(!contactsExpanded)}
            variant="secondary"
            style={styles.toggleButton}
          />

          {contactsExpanded && (
            <View style={styles.contactsList}>
              {partner.contacts.map((contact) => (
                <View key={contact.id} style={styles.contactItem}>
                  <View style={styles.contactHeader}>
                    <Text style={styles.contactName}>{contact.name}</Text>
                    <Badge label={contact.role} color="#666" />
                  </View>
                  <Text style={styles.contactDetail}>Tel: {contact.phone}</Text>
                  <Text style={styles.contactDetail}>Email: {contact.email}</Text>
                </View>
              ))}
            </View>
          )}
        </View>
      )}
    </Card>
  );
});

const styles = StyleSheet.create({
  info: {
    gap: 8,
  },
  row: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 4,
    borderBottomWidth: 1,
    borderBottomColor: '#eee',
  },
  label: {
    fontFamily: 'monospace',
    fontSize: 12,
    textTransform: 'uppercase',
    color: '#666',
  },
  value: {
    fontFamily: 'monospace',
    fontSize: 12,
    fontWeight: 'bold',
  },
  contactsSection: {
    marginTop: 16,
  },
  toggleButton: {
    marginBottom: 12,
  },
  contactsList: {
    gap: 12,
  },
  contactItem: {
    padding: 12,
    backgroundColor: '#f5f5f5',
    borderWidth: 1,
    borderColor: '#ddd',
  },
  contactHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  contactName: {
    fontFamily: 'monospace',
    fontSize: 14,
    fontWeight: 'bold',
  },
  contactDetail: {
    fontFamily: 'monospace',
    fontSize: 11,
    color: '#666',
    marginTop: 4,
  },
});
